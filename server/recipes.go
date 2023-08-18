package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"recipes/models"
	storemodels "recipes/store/models"

	"github.com/gin-gonic/gin"
)

var (
	GetRecipesTimeout = time.Duration(5 * time.Second)
)

// GetRecipes returns recipes
func (s *Server) GetRecipes(c *gin.Context) {
	query := c.Query("query")
	recipeEventID := c.Query("recipe_event_id")

	var ingredientID uint64
	if ingredientIDstr := c.Query("ingredient_id"); ingredientIDstr != "" {
		var err error
		ingredientID, err = strconv.ParseUint(c.Query("ingredient_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var tagID uint64
	if tagIDstr := c.Query("tag_id"); tagIDstr != "" {
		var err error
		tagID, err = strconv.ParseUint(c.Query("tag_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipesTimeout)
	defer cancel()
	recipes, err := s.getRecipes(ctx, query, recipeEventID, ingredientID, tagID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, recipes)
}

func (s *Server) getRecipes(ctx context.Context, query, recipeEventID string, ingredientID, tagID uint64) ([]*models.Recipe, error) {
	if recipeEventID != "" {
		return s.getRecipesByRecipeEventID(ctx, recipeEventID)
	}
	if ingredientID != 0 {
		return s.getRecipesByIngredient(ctx, ingredientID)
	}
	if tagID != 0 {
		return s.getRecipesByTag(ctx, tagID)
	}
	if query != "" {
		return s.getRecipesByQuery(ctx, query)
	}
	return nil, nil
}

func (s *Server) getRecipesByQuery(ctx context.Context, query string) ([]*models.Recipe, error) {
	// TODO: implement
	return nil, nil
}

func (s *Server) getRecipesByRecipeEventID(ctx context.Context, recipeEventID string) ([]*models.Recipe, error) {
	storeRecipes, err := s.store.GetRecipeEventRecipes(ctx, recipeEventID)
	if err != nil {
		return nil, err
	}

	return s.storeRecipesToAPI(ctx, storeRecipes)
}

func (s *Server) getRecipesByIngredient(ctx context.Context, ingredientID uint64) ([]*models.Recipe, error) {
	if _, err := s.store.ExistsIngredient(ctx, ingredientID); err != nil {
		return nil, err
	}
	storeRecipes, err := s.store.GetIngredientRecipes(ctx, ingredientID)
	if err != nil {
		return nil, err
	}

	return s.storeRecipesToAPI(ctx, storeRecipes)
}

func (s *Server) getRecipesByTag(ctx context.Context, tagID uint64) ([]*models.Recipe, error) {
	if _, err := s.store.ExistsTag(ctx, tagID); err != nil {
		return nil, err
	}
	storeRecipes, err := s.store.GetTagRecipes(ctx, tagID)
	if err != nil {
		return nil, err
	}

	return s.storeRecipesToAPI(ctx, storeRecipes)
}

func (s *Server) storeRecipesToAPI(ctx context.Context, storeRecipes []*storemodels.Recipe) ([]*models.Recipe, error) {
	recipes := make([]*models.Recipe, 0, len(storeRecipes))
	for _, storeRecipe := range storeRecipes {
		storeIngredients, storeQuant, err := s.store.GetRecipeIngredients(ctx, storeRecipe.ID)
		if err != nil {
			return nil, err
		}
		storeTags, err := s.store.GetRecipeTags(ctx, storeRecipe.ID)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, storeRecipeToAPI(storeRecipe, storeIngredients, storeQuant, storeTags))
	}
	return recipes, nil
}

func storeRecipeToAPI(storeRecipe *storemodels.Recipe, storeIngredient []*storemodels.Ingredient, storeQuant []*storemodels.RecipeToIngredient, storeTags []*storemodels.Tag) *models.Recipe {
	recipe := &models.Recipe{
		ID:      storeRecipe.ID,
		Name:    storeRecipe.Name,
		Variant: storeRecipe.Variant,
	}

	ingredients := make([]*models.Ingredient, 0, len(storeIngredient))
	for i, storeIngredient := range storeIngredient {
		ingredients = append(ingredients, storeIngredientToAPI(storeIngredient, storeQuant[i]))
	}

	tags := make([]*models.Tag, 0, len(storeTags))
	for _, storeTag := range storeTags {
		tags = append(tags, storeTagToAPI(storeTag))
	}

	recipe.Tags = tags
	recipe.Ingredients = ingredients

	return recipe
}
