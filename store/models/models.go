package models

type RecipeEvent struct {
	Id           string `db:"id"`
	ScheduleDate int64  `db:"schedule_date"`
	Title        string `db:"title"`
	Description  string `db:"description"`
}

type RecipeEventToRecipe struct {
	RecipeEventId string `db:"recipe_event_id"`
	RecipeId      string `db:"recipe_id"`
}

type Recipe struct {
	Id        string `db:"id"`
	Name      string `db:"name"`
	Variant   string `db:"variant"`
	CreatedOn int64  `db:"created_on"`
}

type RecipeTag struct {
	RecipeId string `db:"recipe_id"`
	Tag      string `db:"tag"`
}

type Ingredient struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

type RecipeIngredient struct {
	RecipeId     string `db:"recipe_id"`
	IngredientId string `db:"ingredient_id"`
	Quantity     int    `db:"quantity"`
	Unit         string `db:"unit"`
	Size         string `db:"size"`
}
