package models

type Ingredient struct {
	ID       uint64
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
