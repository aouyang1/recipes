package store

import (
	"context"
	"errors"

	"recipes/store/models"
)

var (
	ErrRecipeToIngredientNotFound = errors.New("recipe to ingredient not found")
)

func (c *Client) UpsertRecipeToIngredient(ctx context.Context, r2i *models.RecipeToIngredient) error {
	if r2i == nil {
		return nil
	}
	if r2i.RecipeID == 0 {
		return ErrInvalidRecipe
	}
	if r2i.IngredientID == 0 {
		return ErrInvalidIngredient
	}

	if _, err := c.ExistsRecipe(ctx, r2i.RecipeID); err != nil {
		return err
	}

	if _, err := c.ExistsIngredient(ctx, r2i.IngredientID); err != nil {
		return err
	}

	_, err := c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_to_ingredient (recipe_id, ingredient_id, quantity, unit, size)
			  VALUES (:recipe_id, :ingredient_id, :quantity, :unit, :size)
         ON DUPLICATE KEY UPDATE quantity = :quantity, unit = :unit, size = :size`,
		r2i,
	)
	return err
}

func (c *Client) GetRecipeIngredients(ctx context.Context, recipeID uint64) ([]*models.Ingredient, []*models.RecipeToIngredient, error) {
	if recipeID == 0 {
		return nil, nil, ErrInvalidRecipe
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name, quantity, unit, size
		   FROM ingredient
		   JOIN (SELECT ingredient_id, quantity, unit, size
		           FROM recipe_to_ingredient
		          WHERE recipe_id = ?) as r2i
			 ON ingredient.id = r2i.ingredient_id`,
		recipeID,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var res []*models.Ingredient
	var quant []*models.RecipeToIngredient
	for rows.Next() {
		var id uint64
		var name string
		var quantity float64
		var unit string
		var size string
		if err := rows.Scan(&id, &name, &quantity, &unit, &size); err != nil {
			return nil, nil, err
		}
		res = append(res, &models.Ingredient{
			ID:   id,
			Name: name,
		})
		quant = append(quant, &models.RecipeToIngredient{
			Quantity: quantity,
			Unit:     unit,
			Size:     size,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return res, quant, nil
}

func (c *Client) GetIngredientRecipes(ctx context.Context, ingredientID uint64) ([]*models.Recipe, error) {
	if ingredientID == 0 {
		return nil, nil
	}

	rows, err := c.conn.QueryxContext(
		ctx,
		`SELECT id, name, variant, created_on
		   FROM recipe 
		   JOIN (SELECT recipe_id
		           FROM recipe_to_ingredient
				  WHERE ingredient_id = ?) as r2i 
			 ON recipe.id = r2i.recipe_id`,
		ingredientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.StructScan(&recipe); err != nil {
			return nil, err
		}
		res = append(res, &recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

type DeleteRecipeToIngredientParams struct {
	RecipeID     uint64 `db:"recipe_id"`
	IngredientID uint64 `db:"ingredient_id"`
}

func (c *Client) DeleteRecipeToIngredient(ctx context.Context, params *DeleteRecipeToIngredientParams) error {
	if params == nil {
		return nil
	}
	if params.RecipeID == 0 {
		return ErrInvalidRecipe
	}
	if params.IngredientID == 0 {
		return ErrInvalidIngredient
	}

	result, err := c.conn.NamedExecContext(
		ctx,
		`DELETE FROM recipe_to_ingredient
		       WHERE recipe_id = :recipe_id
			     AND ingredient_id = :ingredient_id`,
		params,
	)
	if err != nil {
		return err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return ErrRecipeToIngredientNotFound
	}
	return nil
}
