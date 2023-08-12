package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"recipes/models"
	"recipes/store"
	storemodels "recipes/store/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	GOOGLE_API_CREDENTIALS_PATH = os.Getenv("GOOGLE_API_CREDENTIALS_PATH")

	ErrCalendarNotFound = errors.New("calendar not found")

	TimeoutUpsertRecipeEvent = time.Duration(3 * time.Second)
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := path.Join(GOOGLE_API_CREDENTIALS_PATH, "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()

	b, err := os.ReadFile(path.Join(GOOGLE_API_CREDENTIALS_PATH, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calendarSummary := "Dinner Plans"
	calId, err := GetCalendarId(srv, calendarSummary)
	if err != nil {
		log.Fatalf("error while fetching for calendar, %q, %v", calendarSummary, err)
	}
	calEvents, err := GetAllEvents(srv, calId)
	if err != nil {
		log.Fatalf("error while fetching events for calendar, %q, %v", calendarSummary, err)
	}

	recEvents, err := ToRecipeEvents(calEvents)
	if err != nil {
		log.Fatalf("failed to parse all events into recipe events, %v", err)
	}

	storeClient, err := store.NewClient(
		&mysql.Config{
			User:   os.Getenv("USER_MYSQL_USERNAME"),
			Passwd: os.Getenv("USER_MYSQL_PASSWORD"),
			DBName: os.Getenv("USER_MYSQL_DB"),
			Net:    "tcp",
			Addr:   "localhost",
		},
	)
	if err != nil {
		log.Fatalf("failed to create store client, %v", err)
	}

	for _, e := range recEvents {
		// don't track events that have no description
		if e.Description == "" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), TimeoutUpsertRecipeEvent)
		storeRecipeEvent := &storemodels.RecipeEvent{
			ID:           e.ID,
			ScheduleDate: e.Date.Unix(),
			Title:        e.Title,
			Description:  e.Description,
		}
		if err := storeClient.UpsertRecipeEventContext(ctx, storeRecipeEvent); err != nil {
			log.Fatalf("failed to upsert recipe event, %v, %v", e, err)
		}
		cancel()
		fmt.Println(e.Title)
	}
}

func GetCalendarId(srv *calendar.Service, calSummary string) (string, error) {
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve events: %w", err)
	}

	// expecting only one match
	for _, item := range calendars.Items {
		if item.Summary == calSummary {
			return item.Id, nil
		}
	}
	return "", ErrCalendarNotFound
}

func GetAllEvents(srv *calendar.Service, calId string) ([]*calendar.Event, error) {
	return getEvents(srv, calId, "")
}

func getEvents(srv *calendar.Service, calId, pageToken string) ([]*calendar.Event, error) {
	events, err := srv.Events.List(calId).
		ShowDeleted(false).
		SingleEvents(true).
		PageToken(pageToken).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events: %w", err)
	}

	res := events.Items
	if events.NextPageToken != "" {
		nextRes, err := getEvents(srv, calId, events.NextPageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch next page of events for token, %s, %w", events.NextPageToken, err)
		}
		res = append(res, nextRes...)
	}
	return res, nil
}

func ToRecipeEvents(calEvents []*calendar.Event) ([]*models.RecipeEvent, error) {
	recEvents := make([]*models.RecipeEvent, 0, len(calEvents))
	for _, e := range calEvents {
		date := e.Start.Date
		if date == "" {
			continue
		}

		dt, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}

		rec := &models.RecipeEvent{
			ID:          e.Id,
			Date:        dt,
			Title:       e.Summary,
			Description: e.Description,
		}
		recEvents = append(recEvents, rec)
	}
	return recEvents, nil
}
