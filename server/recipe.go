package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"recipes/models"
	"recipes/store"
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
	r, err := s.postRecipe(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, r)
}

func (s *Server) postRecipe(ctx context.Context, req *PostRecipeRequest) (*models.Recipe, error) {
	if req == nil || req.Recipe == nil {
		return nil, nil
	}

	if req.Recipe.Name == "" && req.Recipe.Variant == "" {
		return nil, nil
	}

	r := &storemodels.Recipe{
		ID:        req.Recipe.ID(),
		Name:      req.Recipe.Name,
		Variant:   req.Recipe.Variant,
		CreatedOn: time.Now().UTC().Unix(),
	}

	if _, err := s.store.ExistsRecipe(ctx, r.Name, r.Variant); err != nil {
		if !errors.Is(err, store.ErrRecipeNotFound) {
			return nil, err
		}
		if insertErr := s.store.InsertRecipe(ctx, r); insertErr != nil {
			return nil, err
		}
	}

	if err := s.store.UpsertRecipeEventToRecipe(ctx, req.EventID, req.Recipe.Name, req.Recipe.Variant); err != nil {
		return nil, err
	}
	return req.Recipe, nil
}
