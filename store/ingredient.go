package store

import (
	"context"
	"database/sql"
	"errors"
	"recipes/store/models"
)

var (
	ErrInvalidIngredient        = errors.New("empty ingredient name")
	ErrDuplicateIngredient      = errors.New("cannot insert duplicate ingredient")
	ErrIngredientNotFound       = errors.New("ingredient not found")
	ErrIngredientInUseByRecipes = errors.New("ingredient is in use by recipes")
)

func (c *Client) UpsertIngredient(ctx context.Context, ingredient *models.Ingredient) (uint64, error) {
	iid, err := c.ExistsIngredientName(ctx, ingredient.Name)
	if err != nil {
		if !errors.Is(err, ErrIngredientNotFound) {
			return 0, err
		}
	}

	// If the ingredient name exists, return early. we don't want too update the id
	// of the ingredient.
	if iid != 0 && ingredient.ID != iid {
		return iid, ErrDuplicateIngredient
	}

	if ingredient.ID == 0 {
		result, err := c.conn.NamedExecContext(
			ctx,
			`INSERT INTO ingredient (name)
			  VALUES (:name)`,
			ingredient,
		)
		if err != nil {
			return 0, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		return uint64(id), nil
	}

	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO ingredient (id, name)
			  VALUES (:id, :name)
		 ON DUPLICATE KEY UPDATE name = :name`,
		ingredient,
	)
	if err != nil {
		return 0, err
	}

	return ingredient.ID, nil
}

func (c *Client) ExistsIngredient(ctx context.Context, ingredientID uint64) (uint64, error) {
	if ingredientID == 0 {
		return 0, ErrInvalidIngredient
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM ingredient WHERE id = ?`,
		ingredientID,
	)
	if err := row.Scan(&ingredientID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrIngredientNotFound
		}
		return 0, err
	}
	return ingredientID, nil
}

func (c *Client) ExistsIngredientName(ctx context.Context, ingredient string) (uint64, error) {
	if ingredient == "" {
		return 0, nil
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM ingredient WHERE name = ?`,
		ingredient,
	)
	var ingredientID uint64
	if err := row.Scan(&ingredientID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrIngredientNotFound
		}
		return 0, err
	}
	return ingredientID, nil
}

func (c *Client) GetIngredients(ctx context.Context) ([]*models.Ingredient, error) {
	rows, err := c.conn.QueryxContext(
		ctx,
		`SELECT id, name, count(recipe_to_ingredient.recipe_id) as cnt
		   FROM ingredient
	  LEFT JOIN recipe_to_ingredient
		     ON ingredient.id = recipe_to_ingredient.ingredient_id
	   GROUP BY ingredient.id
	   ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*models.Ingredient
	for rows.Next() {
		var ingredient models.Ingredient
		if err := rows.StructScan(&ingredient); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ingredient)
	}
	return ingredients, nil
}

func (c *Client) DeleteIngredient(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidIngredient
	}

	recipes, err := c.GetIngredientRecipes(ctx, id)
	if err != nil {
		return err
	}

	if len(recipes) > 0 {
		return ErrIngredientInUseByRecipes
	}

	_, err = c.conn.ExecContext(
		ctx,
		`DELETE FROM ingredient WHERE id = ?`,
		id,
	)
	return err
}
