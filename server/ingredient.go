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
	PostIngredientTimeout   = time.Duration(5 * time.Second)
	DeleteIngredientTimeout = time.Duration(5 * time.Second)
)

// PostIngredient adds an ingredient to the db
func (s *Server) PostIngredient(c *gin.Context) {
	var req models.Ingredient
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostIngredientTimeout)
	defer cancel()
	ingredientID, err := s.postIngredient(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	req.ID = ingredientID
	c.JSON(http.StatusCreated, req)
}

func (s *Server) postIngredient(ctx context.Context, ingredient *models.Ingredient) (uint64, error) {
	i := &storemodels.Ingredient{
		ID:   ingredient.ID,
		Name: ingredient.Name,
	}
	return s.store.UpsertIngredient(ctx, i)
}

// DeleteIngredient removes an ingredient from the db
func (s *Server) DeleteIngredient(c *gin.Context) {
	var req models.Ingredient
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeleteIngredientTimeout)
	defer cancel()
	if err := s.deleteIngredient(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) deleteIngredient(ctx context.Context, ingredient *models.Ingredient) error {
	return s.store.DeleteIngredient(ctx, ingredient.ID)
}
