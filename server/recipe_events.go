package main

import (
	"context"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetRecipeEventsTimeout = time.Duration(5 * time.Second)
)

// GetRecipeEvents returns recipe events
func (s *Server) GetRecipeEvents(c *gin.Context) {
	recipeName := c.Query("recipe_name")
	recipeVariant := c.Query("recipe_variant")

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipeEventsTimeout)
	defer cancel()
	recipeEvents, err := s.getRecipeEvents(ctx, recipeName, recipeVariant)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, recipeEvents)
}

func (s *Server) getRecipeEvents(ctx context.Context, name, variant string) ([]*models.RecipeEvent, error) {
	if name == "" && variant == "" {
		return s.getAllRecipeEvents(ctx)
	}
	return s.getRecipeEventsByRecipe(ctx, name, variant)
}

func (s *Server) getAllRecipeEvents(ctx context.Context) ([]*models.RecipeEvent, error) {
	storeRecipeEvents, err := s.store.GetRecipeEvents(ctx)
	if err != nil {
		return nil, err
	}

	return storeRecipeEventsToAPI(storeRecipeEvents)
}

func (s *Server) getRecipeEventsByRecipe(ctx context.Context, name, variant string) ([]*models.RecipeEvent, error) {
	storeRecipeEvents, err := s.store.GetRecipeRecipeEvents(ctx, name, variant)
	if err != nil {
		return nil, err
	}

	return storeRecipeEventsToAPI(storeRecipeEvents)
}

func storeRecipeEventsToAPI(storeRecipeEvents []*storemodels.RecipeEvent) ([]*models.RecipeEvent, error) {
	recipeEvents := make([]*models.RecipeEvent, 0, len(storeRecipeEvents))
	for _, storeRecipeEvent := range storeRecipeEvents {
		re, err := models.NewRecipeEvent(
			storeRecipeEvent.ID,
			time.Unix(storeRecipeEvent.ScheduleDate, 0).UTC(),
			storeRecipeEvent.Title,
			storeRecipeEvent.Description,
		)
		if err != nil {
			return nil, err
		}
		re.Count = storeRecipeEvent.Count
		recipeEvents = append(recipeEvents, re)
	}
	return recipeEvents, nil
}
