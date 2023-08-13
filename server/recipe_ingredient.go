package main

import (
	"context"
	"net/http"
	"time"

	"recipes/store"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	PostRecipeIngredientTimeout   = time.Duration(5 * time.Second)
	DeleteRecipeIngredientTimeout = time.Duration(5 * time.Second)
)

type RecipeIngredientRequest struct {
	RecipeName    string `json:"recipe_name"`
	RecipeVariant string `json:"recipe_variant"`
	Ingredient    string `json:"tag"`
	Quantity      int    `json:"quantity"`
	Unit          string `json:"unit"`
	Size          string `json:"size"`
}

// PostRecipeIngredient adds an ingredient to a recipe
func (s *Server) PostRecipeIngredient(c *gin.Context) {
	var req RecipeIngredientRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostRecipeIngredientTimeout)
	defer cancel()
	if err := s.postRecipeIngredient(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) postRecipeIngredient(ctx context.Context, req *RecipeIngredientRequest) error {
	if req == nil {
		return nil
	}
	r2i := storemodels.RecipeToIngredient{
		Quantity: req.Quantity,
		Unit:     req.Unit,
		Size:     req.Size,
	}
	return s.store.UpsertRecipeToIngredient(ctx, req.RecipeName, req.RecipeVariant, req.Ingredient, r2i)
}

// DeleteRecipeIngredient removes an ingredient from a recipe
func (s *Server) DeleteRecipeIngredient(c *gin.Context) {
	var req RecipeIngredientRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeleteRecipeIngredientTimeout)
	if err := s.deleteRecipeIngredient(ctx, &req); err != nil {
		cancel()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cancel()
	c.Status(http.StatusOK)
}

func (s *Server) deleteRecipeIngredient(ctx context.Context, req *RecipeIngredientRequest) error {
	if req == nil {
		return nil
	}
	params := store.DeleteRecipeToIngredientParams{
		RecipeName:    req.RecipeName,
		RecipeVariant: req.RecipeVariant,
		Ingredient:    req.Ingredient,
	}
	return s.store.DeleteRecipeToIngredient(ctx, params)
}
