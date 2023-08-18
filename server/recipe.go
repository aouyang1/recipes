package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetRecipeTimeout = time.Duration(5 * time.Second)
)

type PostRecipeRequest struct {
	EventID string         `json:"recipe_event_id"`
	Recipe  *models.Recipe `json:"recipe"`
}

// PostRecipe adds a recipe and links it to an event
func (s *Server) PostRecipe(c *gin.Context) {
	var req PostRecipeRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipeTimeout)
	defer cancel()
	if err := s.postRecipe(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (s *Server) postRecipe(ctx context.Context, req *PostRecipeRequest) error {
	if req == nil || req.Recipe == nil {
		return nil
	}

	if req.Recipe.Name == "" && req.Recipe.Variant == "" {
		return nil
	}

	r := &storemodels.Recipe{
		Name:      req.Recipe.Name,
		Variant:   req.Recipe.Variant,
		CreatedOn: time.Now().UTC().Unix(),
	}

	recipeID, err := s.store.UpsertRecipe(ctx, r)
	if err != nil {
		return err
	}

	return s.store.UpsertRecipeEventToRecipe(ctx, req.EventID, recipeID)
}

// UpdateRecipe updates a recipe
func (s *Server) UpdateRecipe(c *gin.Context) {
	var req models.Recipe
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipeTimeout)
	defer cancel()
	if err := s.updateRecipe(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) updateRecipe(ctx context.Context, req *models.Recipe) error {
	r := &storemodels.Recipe{
		ID:        req.ID,
		Name:      req.Name,
		Variant:   req.Variant,
		CreatedOn: time.Now().UTC().Unix(),
	}

	_, err := s.store.UpsertRecipe(ctx, r)
	return err
}
