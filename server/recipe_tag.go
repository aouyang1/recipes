package main

import (
	"context"
	"net/http"
	"time"

	"recipes/store"

	"github.com/gin-gonic/gin"
)

var (
	PostRecipeTagTimeout   = time.Duration(5 * time.Second)
	DeleteRecipeTagTimeout = time.Duration(5 * time.Second)
)

type RecipeTagRequest struct {
	RecipeName    string `json:"recipe_name"`
	RecipeVariant string `json:"recipe_variant"`
	Tag           string `json:"tag"`
}

// PostRecipeTag adds a tag to a recipe
func (s *Server) PostRecipeTag(c *gin.Context) {
	var req RecipeTagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), PostRecipeTagTimeout)
	defer cancel()
	if err := s.postRecipeTag(ctx, &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) postRecipeTag(ctx context.Context, req *RecipeTagRequest) error {
	if req == nil {
		return nil
	}
	return s.store.UpsertRecipeToTag(ctx, req.RecipeName, req.RecipeVariant, req.Tag)
}

// DeleteRecipeTag removes a tag from a recipe
func (s *Server) DeleteRecipeTag(c *gin.Context) {
	var req RecipeTagRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeleteRecipeTagTimeout)
	if err := s.deleteRecipeTag(ctx, &req); err != nil {
		cancel()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cancel()
	c.Status(http.StatusOK)
}

func (s *Server) deleteRecipeTag(ctx context.Context, req *RecipeTagRequest) error {
	if req == nil {
		return nil
	}
	return s.store.DeleteRecipeToTag(ctx, store.DeleteRecipeToTagParams(*req))
}
