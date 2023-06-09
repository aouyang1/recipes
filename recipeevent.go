package main

import (
	"fmt"
	"time"
)

// RecipeEvent represents an occurrence of a recipe taken from the calendar
type RecipeEvent struct {
	Id          string // max 60 character id
	Date        time.Time
	Title       string
	Description string
}

func (r RecipeEvent) String() string {
	return fmt.Sprintf("id: %s, date: %s, title: %s", r.Id, r.Date.Format("2006-01-02"), r.Title)
}

// Recipe represents a normalized representation of a recipe
type Recipe struct {
	Id          string
	Name        string   // Chicken Marsala
	Variant     string   // bonapetit, new york times, etc.
	Tags        []string // Italian, Noodles, Meat, etc.
	Ingredients []*Ingredient
}

type Ingredient struct {
	Id       string
	Name     string
	Quantity int
	Size     Size
	Unit     Unit
}

type Size string

const (
	Size_Large  = "lg"
	Size_Medium = "md"
	Size_Small  = "sm"
)

type Unit string

const (
	Unit_Ounce      = "oz."
	Unit_Pound      = "lb."
	Unit_Cup        = "cup"
	Unit_Teaspoon   = "tsp."
	Unit_Tablespoon = "tbsp."
)
