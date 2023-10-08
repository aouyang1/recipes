package models

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"
)

var (
	ErrNoRecipeEventID    = errors.New("no recipe event id")
	ErrNoRecipeEventDate  = errors.New("no recipe event date")
	ErrNoRecipeEventTitle = errors.New("no recipe event title")
)

// RecipeEvent represents an occurrence of a recipe taken from the calendar
type RecipeEvent struct {
	ID          string    `json:"id"` // max 60 character id
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URLs        []string  `json:"url_links"`
	Count       int       `json:"count"`
}

func (r RecipeEvent) String() string {
	return fmt.Sprintf("id: %s, date: %s, title: %s", r.ID, r.Date.Format("2006-01-02"), r.Title)
}

func NewRecipeEvent(id string, date time.Time, title, description string) (*RecipeEvent, error) {
	if id == "" {
		return nil, ErrNoRecipeEventID
	}
	if date.IsZero() {
		return nil, ErrNoRecipeEventDate
	}
	if title == "" {
		return nil, ErrNoRecipeEventTitle
	}

	rxStrict := xurls.Strict()
	urls := rxStrict.FindAllString(description, -1)
	distinctURLs := make(map[string]struct{})
	for _, u := range urls {
		distinctURLs[cleanURL(u)] = struct{}{}
	}
	urls = make([]string, 0, len(distinctURLs))
	for u := range distinctURLs {
		urls = append(urls, u)
	}
	return &RecipeEvent{
		ID:          id,
		Date:        date,
		Title:       title,
		URLs:        urls,
		Description: description,
	}, nil
}

func cleanURL(urlStr string) string {
	if strings.Contains(urlStr, "www.google.com") {
		//"https://www.google.com/url?q=https://carnaldish.com/recipes/the-best-homemade-smash-burgers/&sa=D&source=calendar&usd=2&usg=AOvVaw2Lz8Qatlh7ro-wBOxTj1Bn"
		urlParsed, _ := url.Parse(urlStr)
		query := urlParsed.Query().Get("q")
		if len(query) > 0 {
			return query
		}
	}
	return urlStr
}
