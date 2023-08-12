package models

type RecipeEvent struct {
	ID           string `db:"id"`
	ScheduleDate int64  `db:"schedule_date"`
	Title        string `db:"title"`
	Description  string `db:"description"`
}

type RecipeEventToRecipe struct {
	RecipeEventID string `db:"recipe_event_id"`
	RecipeID      uint64 `db:"recipe_id"`
}

type Recipe struct {
	ID        uint64 `db:"id"`
	Name      string `db:"name"`
	Variant   string `db:"variant"`
	CreatedOn int64  `db:"created_on"`
}

type Tag struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

type RecipeToTag struct {
	RecipeID uint64 `db:"recipe_id"`
	TagID    uint64 `db:"tag_id"`
}

type Ingredient struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

type RecipeToIngredient struct {
	RecipeID     uint64 `db:"recipe_id"`
	IngredientID uint64 `db:"ingredient_id"`

	Quantity int    `db:"quantity"`
	Unit     string `db:"unit"`
	Size     string `db:"size"`
}
