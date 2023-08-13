package main

import (
	"context"
	"net/http"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	PostIngredientTimeout = time.Duration(5 * time.Second)
)

func (s *Server) PostIngredient(c *gin.Context) {
	var req models.Ingredient
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostIngredientTimeout)
	defer cancel()
	if err := s.postIngredient(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) postIngredient(ctx context.Context, ingredient *models.Ingredient) error {
	i := &storemodels.Ingredient{
		ID:   ingredient.ID(),
		Name: ingredient.Name,
	}
	return s.store.InsertIngredient(ctx, i)
}
