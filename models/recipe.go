package models

import (
	"github.com/cespare/xxhash"
)

// Recipe represents a normalized representation of a recipe
type Recipe struct {
	Name    string // Chicken Marsala
	Variant string // bonapetit, new york times, etc.

	Tags        map[Tag]struct{}       // Italian, Noodles, Meat, etc.
	Ingredients map[string]*Ingredient // key is name of ingredient
}

func (r *Recipe) ID() uint64 {
	return xxhash.Sum64String(r.Name + "/" + r.Variant)
}
