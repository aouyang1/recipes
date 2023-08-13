package main

import (
	"context"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetRecipeTimeout = time.Duration(5 * time.Second)
)

// PostRecipe adds a recipe
func (s *Server) PostRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipeTimeout)
	r, err := s.postRecipe(ctx, &recipe)
	if err != nil {
		cancel()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cancel()
	c.JSON(200, r)
}

func (s *Server) postRecipe(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	if recipe == nil {
		return nil, nil
	}
	if recipe.Name == "" && recipe.Variant == "" {
		return nil, nil
	}

	r := &storemodels.Recipe{
		ID:        recipe.ID(),
		Name:      recipe.Name,
		Variant:   recipe.Variant,
		CreatedOn: time.Now().UTC().Unix(),
	}
	return recipe, s.store.InsertRecipe(ctx, r)
}
