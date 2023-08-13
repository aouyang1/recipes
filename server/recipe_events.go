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
	recipeEvents, err := s.getRecipeEvents(ctx, recipeName, recipeVariant)
	if err != nil {
		cancel()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cancel()
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

	return storeRecipeEventsToAPI(storeRecipeEvents), nil
}

func (s *Server) getRecipeEventsByRecipe(ctx context.Context, name, variant string) ([]*models.RecipeEvent, error) {
	storeRecipeEvents, err := s.store.GetRecipeRecipeEvents(ctx, name, variant)
	if err != nil {
		return nil, err
	}

	return storeRecipeEventsToAPI(storeRecipeEvents), nil
}

func storeRecipeEventsToAPI(storeRecipeEvents []*storemodels.RecipeEvent) []*models.RecipeEvent {
	recipeEvents := make([]*models.RecipeEvent, 0, len(storeRecipeEvents))
	for _, storeRecipeEvent := range storeRecipeEvents {
		recipeEvents = append(recipeEvents, &models.RecipeEvent{
			ID:          storeRecipeEvent.ID,
			Date:        time.Unix(storeRecipeEvent.ScheduleDate, 0).UTC(),
			Title:       storeRecipeEvent.Title,
			Description: storeRecipeEvent.Description,
		})
	}
	return recipeEvents
}
