package models

import (
	"github.com/cespare/xxhash"
)

// Recipe represents a normalized representation of a recipe
type Recipe struct {
	Name    string `json:"name"`    // Chicken Marsala
	Variant string `json:"variant"` // url link or bonapetit, new york times, etc.

	Tags        []Tag         `json:"tags"` // Italian, Noodles, Meat, etc.
	Ingredients []*Ingredient `json:"ingredients"`
}

func (r *Recipe) ID() uint64 {
	return xxhash.Sum64String(r.Name + "/" + r.Variant)
}
