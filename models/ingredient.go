package models

type Ingredient struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	Quantity int  `json:"quantity"`
	Size     Size `json:"size,omitempty"`
	Unit     Unit `json:"unit,omitempty"`
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
