package models

type Ingredient struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	RecipeCount int    `json:"recipe_count"`

	Quantity float64 `json:"quantity"`
	Size     Size    `json:"size,omitempty"`
	Unit     Unit    `json:"unit,omitempty"`
}

type Size string

const (
	SizeLarge  = "lg"
	SizeMedium = "md"
	SizeSmall  = "sm"
)

type Unit string

const (
	UnitMilliliter = "ml."
	UnitGram       = "g."
	UnitOunce      = "oz."
	UnitPound      = "lb."
	UnitCup        = "cup"
	UnitTeaspoon   = "tsp."
	UnitTablespoon = "tbsp."
)
