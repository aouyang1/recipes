package models

import (
	"github.com/cespare/xxhash"
)

type Ingredient struct {
	Name string `json:"name"`

	Quantity int  `json:"quantity"`
	Size     Size `json:"size,omitempty"`
	Unit     Unit `json:"unit,omitempty"`
}

func (i *Ingredient) ID() uint64 {
	return xxhash.Sum64String(i.Name)
}

type Size string

const (
	SizeLarge  = "lg"
	SizeMedium = "md"
	SizeSmall  = "sm"
)

type Unit string

const (
	UnitOunce      = "oz."
	UnitPound      = "lb."
	UnitCup        = "cup"
	UnitTeaspoon   = "tsp."
	UnitTablespoon = "tbsp."
)
