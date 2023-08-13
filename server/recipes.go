package main

import (
	"context"
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
	// query params in order of priority
	query := c.Query("query")
	recipeEventID := c.Query("recipe_event_id")
	ingredient := c.Query("ingredient")
	tag := c.Query("tag")

	ctx, cancel := context.WithTimeout(context.Background(), GetRecipesTimeout)
	recipes, err := s.getRecipes(ctx, query, recipeEventID, ingredient, tag)
	if err != nil {
		cancel()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cancel()
	c.JSON(200, recipes)
}

func (s *Server) getRecipes(ctx context.Context, query, recipeEventID, ingredient, tag string) ([]*models.Recipe, error) {
	if query != "" {
		return s.getRecipesByQuery(ctx, query)
	}
	if recipeEventID != "" {
		return s.getRecipesByRecipeEventID(ctx, recipeEventID)
	}
	if ingredient != "" {
		return s.getRecipesByIngredient(ctx, ingredient)
	}
	if tag != "" {
		return s.getRecipesByTag(ctx, tag)
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

func (s *Server) getRecipesByIngredient(ctx context.Context, ingredient string) ([]*models.Recipe, error) {
	storeRecipes, err := s.store.GetIngredientRecipes(ctx, []string{ingredient})
	if err != nil {
		return nil, err
	}

	return s.storeRecipesToAPI(ctx, storeRecipes)
}

func (s *Server) getRecipesByTag(ctx context.Context, tag string) ([]*models.Recipe, error) {
	storeRecipes, err := s.store.GetTagRecipes(ctx, []string{tag})
	if err != nil {
		return nil, err
	}

	return s.storeRecipesToAPI(ctx, storeRecipes)
}

func (s *Server) storeRecipesToAPI(ctx context.Context, storeRecipes []*storemodels.Recipe) ([]*models.Recipe, error) {
	recipes := make([]*models.Recipe, 0, len(storeRecipes))
	for _, storeRecipe := range storeRecipes {
		storeIngredients, storeQuant, err := s.store.GetRecipeIngredients(ctx, storeRecipe.Name, storeRecipe.Variant)
		if err != nil {
			return nil, err
		}
		storeTags, err := s.store.GetRecipeTags(ctx, storeRecipe.Name, storeRecipe.Variant)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, storeRecipeToAPI(storeRecipe, storeIngredients, storeQuant, storeTags))
	}
	return recipes, nil
}

func storeRecipeToAPI(storeRecipe *storemodels.Recipe, storeIngredient []*storemodels.Ingredient, storeQuant []*storemodels.RecipeToIngredient, storeTags []*storemodels.Tag) *models.Recipe {
	recipe := &models.Recipe{
		Name:    storeRecipe.Name,
		Variant: storeRecipe.Variant,
	}

	ingredients := make([]*models.Ingredient, 0, len(storeIngredient))
	for i, storeIngredient := range storeIngredient {
		ingredients = append(ingredients, storeIngredientToAPI(storeIngredient, storeQuant[i]))
	}

	tags := make([]models.Tag, 0, len(storeTags))
	for _, storeTag := range storeTags {
		tags = append(tags, storeTagToAPI(storeTag))
	}

	recipe.Tags = tags
	recipe.Ingredients = ingredients

	return recipe
}
