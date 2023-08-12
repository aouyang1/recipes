package models

// Recipe represents a normalized representation of a recipe
type Recipe struct {
	ID          uint64
	Name        string // Chicken Marsala
	Variant     string // bonapetit, new york times, etc.
	Tags        []*Tag // Italian, Noodles, Meat, etc.
	Ingredients []*Ingredient
}
