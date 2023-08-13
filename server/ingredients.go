package main

import (
	"recipes/models"
	storemodels "recipes/store/models"
)

func storeIngredientToAPI(storeIngredient *storemodels.Ingredient, storeQuant *storemodels.RecipeToIngredient) *models.Ingredient {
	return &models.Ingredient{
		Name:     storeIngredient.Name,
		Quantity: storeQuant.Quantity,
		Unit:     models.Unit(storeQuant.Unit),
		Size:     models.Size(storeQuant.Size),
	}
}
