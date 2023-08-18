package models

// Recipe represents a normalized representation of a recipe
type Recipe struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`    // Chicken Marsala
	Variant string `json:"variant"` // url link or bonapetit, new york times, etc.

	Tags        []*Tag        `json:"tags"` // Italian, Noodles, Meat, etc.
	Ingredients []*Ingredient `json:"ingredients"`
}
