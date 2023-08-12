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

	NewServiceTimeout        = time.Duration(5 * time.Second)
	UpsertRecipeEventTimeout = time.Duration(3 * time.Second)
)

type CalendarFetcher struct {
	calSvc *calendar.Service
	store  *store.Client
}

func NewCalendarFetcher(calendarSummary string) (*CalendarFetcher, error) {
	b, err := os.ReadFile(path.Join(GOOGLE_API_CREDENTIALS_PATH, "credentials.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client := getClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), NewServiceTimeout)

	svc, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("unable to retrieve Calendar client: %w", err)
	}
	cancel()

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
		return nil, fmt.Errorf("failed to create store client, %w", err)
	}

	cf := &CalendarFetcher{
		calSvc: svc,
		store:  storeClient,
	}
	return cf, nil
}

func (c *CalendarFetcher) Fetch() error {
	calId, err := c.GetCalendarId(calendarSummary)
	if err != nil {
		return fmt.Errorf("error while fetching for calendar, %q, %w", calendarSummary, err)
	}
	calEvents, err := c.GetAllEvents(calId)
	if err != nil {
		return fmt.Errorf("error while fetching events for calendar, %q, %v", calendarSummary, err)
	}

	recEvents, err := ToRecipeEvents(calEvents)
	if err != nil {
		return fmt.Errorf("failed to parse all events into recipe events, %v", err)
	}

	return c.storeEvents(recEvents)
}

func (c *CalendarFetcher) GetCalendarId(calSummary string) (string, error) {
	calendars, err := c.calSvc.CalendarList.List().Do()
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

func (c *CalendarFetcher) GetAllEvents(calId string) ([]*calendar.Event, error) {
	return c.getEvents(calId, "")
}

func (c *CalendarFetcher) getEvents(calId, pageToken string) ([]*calendar.Event, error) {
	events, err := c.calSvc.Events.List(calId).
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
		nextRes, err := c.getEvents(calId, events.NextPageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch next page of events for token, %s, %w", events.NextPageToken, err)
		}
		res = append(res, nextRes...)
	}
	return res, nil
}

func (c *CalendarFetcher) storeEvents(recEvents []*models.RecipeEvent) error {
	for _, e := range recEvents {
		// don't track events that have no description
		if e.Description == "" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), UpsertRecipeEventTimeout)
		storeRecipeEvent := &storemodels.RecipeEvent{
			ID:           e.ID,
			ScheduleDate: e.Date.Unix(),
			Title:        e.Title,
			Description:  e.Description,
		}
		if err := c.store.UpsertRecipeEventContext(ctx, storeRecipeEvent); err != nil {
			cancel()
			return fmt.Errorf("failed to upsert recipe event, %v, %w", e, err)
		}
		cancel()
	}
	return nil
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
