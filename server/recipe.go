package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	recipe, err := s.postRecipe(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

func (s *Server) postRecipe(ctx context.Context, req *PostRecipeRequest) (*models.Recipe, error) {
	if req == nil || req.Recipe == nil {
		return nil, nil
	}

	if req.Recipe.Name == "" && req.Recipe.Variant == "" {
		return nil, nil
	}

	r := &storemodels.Recipe{
		Name:      req.Recipe.Name,
		Variant:   req.Recipe.Variant,
		CreatedOn: time.Now().UTC().Unix(),
	}

	recipeID, err := s.store.UpsertRecipe(ctx, r)
	if err != nil {
		// if already exists then it means we are linking an existing recipe
		if !errors.Is(err, store.ErrDuplicateRecipe) {
			return nil, err
		}
	}

	if err := s.store.UpsertRecipeEventToRecipe(ctx, req.EventID, recipeID); err != nil {
		return nil, err
	}

	recipes, err := s.getRecipes(ctx, "", req.EventID, 0, 0)
	if err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		if recipe.ID == recipeID {
			return recipe, nil
		}
	}

	return nil, store.ErrRecipeNotFound
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

	if _, err := s.store.UpsertRecipe(ctx, r); err != nil {
		return err
	}

	existingTags, err := s.store.GetRecipeTags(ctx, r.ID)
	if err != nil {
		return err
	}

	reqTags := make(map[string]struct{})
	for _, tag := range req.Tags {
		reqTags[tag.Name] = struct{}{}
		if err := s.store.UpsertRecipeToTag(ctx, r.ID, tag.ID); err != nil {
			return err
		}
	}

	// remove tags that are no longer associated with the recipe
	for _, tag := range existingTags {
		if _, ok := reqTags[tag.Name]; !ok {
			params := &store.DeleteRecipeToTagParams{
				RecipeID: r.ID,
				TagID:    tag.ID,
			}
			if err := s.store.DeleteRecipeToTag(ctx, params); err != nil {
				return err
			}
		}
	}

	existingIngredients, _, err := s.store.GetRecipeIngredients(ctx, r.ID)
	if err != nil {
		return err
	}

	reqIngredients := make(map[string]struct{})
	for _, ingredient := range req.Ingredients {
		reqIngredients[ingredient.Name] = struct{}{}
		r2i := &storemodels.RecipeToIngredient{
			RecipeID:     r.ID,
			IngredientID: ingredient.ID,
			Quantity:     ingredient.Quantity,
			Unit:         string(ingredient.Unit),
			Size:         string(ingredient.Size),
		}
		if err := s.store.UpsertRecipeToIngredient(ctx, r2i); err != nil {
			return err
		}
	}

	// remove ingredients that are no longer associated with the recipe
	for _, ingredient := range existingIngredients {
		if _, ok := reqIngredients[ingredient.Name]; !ok {
			params := &store.DeleteRecipeToIngredientParams{
				RecipeID:     r.ID,
				IngredientID: ingredient.ID,
			}
			if err := s.store.DeleteRecipeToIngredient(ctx, params); err != nil {
				return err
			}
		}
	}

	return nil
}
