package store

import "recipes/store/models"

func (c *Client) InsertRecipeEvent(recipeEvent *models.RecipeEvent) error {
	return nil
}

func (c *Client) ExistsRecipeEvent(recipeEvent *models.RecipeEvent) (bool, error) {
	return false, nil
}
