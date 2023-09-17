package main

import (
	"context"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetIngredientsTimeout = time.Duration(5 * time.Second)
)

// GetIngredients returns all ingredients
func (s *Server) GetIngredients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), GetIngredientsTimeout)
	defer cancel()
	ingredients, err := s.getIngredients(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ingredients)
}

func (s *Server) getIngredients(ctx context.Context) ([]*models.Ingredient, error) {
	storeIngredients, err := s.store.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	ingredients := make([]*models.Ingredient, 0, len(storeIngredients))
	for _, storeIngredient := range storeIngredients {
		ingredients = append(ingredients, storeIngredientToAPI(storeIngredient, nil))
	}

	return ingredients, nil
}

func storeIngredientToAPI(storeIngredient *storemodels.Ingredient, storeQuant *storemodels.RecipeToIngredient) *models.Ingredient {
	i := &models.Ingredient{
		ID:          storeIngredient.ID,
		Name:        storeIngredient.Name,
		RecipeCount: storeIngredient.Count,
	}
	if storeQuant != nil {
		i.Quantity = storeQuant.Quantity
		i.Unit = models.Unit(storeQuant.Unit)
		i.Size = models.Size(storeQuant.Size)
	}
	return i
}
