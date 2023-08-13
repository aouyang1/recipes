package store

import (
	"context"
	"errors"
	"strings"

	"recipes/store/models"
)

var (
	ErrRecipeToIngredientNotFound = errors.New("recipe to ingredient not found")
)

func (c *Client) UpsertRecipeToIngredient(ctx context.Context, recipeName, recipeVariant, ingredientName string, r2i models.RecipeToIngredient) error {
	if recipeName == "" {
		return ErrInvalidRecipe
	}
	if ingredientName == "" {
		return ErrInvalidIngredient
	}

	recipeID, err := c.ExistsRecipe(ctx, recipeName, recipeVariant)
	if err != nil {
		return err
	}

	ingredientID, err := c.ExistsIngredient(ctx, ingredientName)
	if err != nil {
		return err
	}

	// override existing recipe id and ingredient id
	r2i.RecipeID = recipeID
	r2i.IngredientID = ingredientID

	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_to_ingredient (recipe_id, ingredient_id, quantity, unit, size)
			  VALUES (:recipe_id, :ingredient_id, :quantity, :unit, :size)
         ON DUPLICATE KEY UPDATE quantity = :quantity, unit = :unit, size = :size`,
		r2i,
	)
	return err
}

func (c *Client) GetRecipeIngredients(ctx context.Context, name, variant string) ([]*models.Ingredient, []*models.RecipeToIngredient, error) {
	if name == "" {
		return nil, nil, ErrInvalidRecipe
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name, quantity, unit, size
		   FROM ingredient
		   JOIN (SELECT ingredient_id, quantity, unit, size
		           FROM recipe_to_ingredient
		          WHERE recipe_id = (SELECT id FROM recipe WHERE name = ? AND variant = ?) as r2i
			 ON ingredient.id = r2i.ingredient_id`,
		strings.ToLower(name),
		LowerVariant(variant),
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
		var quantity int
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

func (c *Client) GetIngredientRecipes(ctx context.Context, ingredients []string) ([]*models.Recipe, error) {
	if len(ingredients) == 0 {
		return nil, nil
	}

	for i, ingredient := range ingredients {
		ingredients[i] = strings.ToLower(ingredient)
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name, variant, created_on
		   FROM recipe 
		   JOIN (SELECT recipe_id
		           FROM recipe_to_ingredient
		           JOIN (SELECT id FROM ingredient WHERE name IN (?))
				     ON recipe_to_ingredient.ingredient_id = ingredient.id) as r2i 
			 ON recipe.id = r2i.recipe_id`,
		ingredients,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Recipe
	for rows.Next() {
		var id uint64
		var name string
		var variant string
		var createdOn int64
		if err := rows.Scan(&id, &name, &variant, &createdOn); err != nil {
			return nil, err
		}
		res = append(res, &models.Recipe{
			ID:        id,
			Name:      name,
			Variant:   variant,
			CreatedOn: createdOn,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

type DeleteRecipeToIngredientParams struct {
	RecipeName    string `db:"recipe_name"`
	RecipeVariant string `db:"recipe_variant"`
	Ingredient    string `db:"ingredient"`
}

func (c *Client) DeleteRecipeToIngredient(ctx context.Context, params DeleteRecipeToIngredientParams) error {
	if params.RecipeName == "" {
		return ErrInvalidRecipe
	}
	if params.Ingredient == "" {
		return ErrInvalidIngredient
	}

	params.RecipeVariant = LowerVariant(params.RecipeVariant)

	result, err := c.conn.NamedExecContext(
		ctx,
		`DELETE FROM recipe_to_ingredient
		       WHERE recipe_id = (SELECT id FROM recipe WHERE name = :recipe_name AND variant = :recipe_variant)
			     AND ingredient_id = (SELECT id FROM ingredient WHERE name = :ingredient)`,
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
