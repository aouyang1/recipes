package models

import (
	"errors"
	"fmt"
	"time"
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
	return &RecipeEvent{
		ID:          id,
		Date:        date,
		Title:       title,
		Description: description,
	}, nil
}
