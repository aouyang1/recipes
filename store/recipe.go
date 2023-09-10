package store

import (
	"context"
	"database/sql"
	"errors"

	"recipes/store/models"
)

var (
	ErrInvalidRecipe   = errors.New("empty recipe name")
	ErrDuplicateRecipe = errors.New("cannot insert duplicate recipe")
	ErrRecipeNotFound  = errors.New("recipe not found")
)

func (c *Client) UpsertRecipe(ctx context.Context, recipe *models.Recipe) (uint64, error) {
	rid, err := c.ExistsRecipeName(ctx, recipe.Name, recipe.Variant)
	if err != nil {
		if !errors.Is(err, ErrRecipeNotFound) {
			return 0, err
		}
	}
	// If the recipe name/variant exists, return early. we don't want too update the id
	// of the recipe.
	if rid != 0 && recipe.ID != rid {
		return rid, ErrDuplicateRecipe
	}

	if recipe.ID == 0 {
		result, err := c.conn.NamedExecContext(
			ctx,
			`INSERT INTO recipe (name, variant, created_on)
			  VALUES (:name, :variant, :created_on)`,
			recipe,
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
		`INSERT INTO recipe (id, name, variant, created_on)
			  VALUES (:id, :name, :variant, :created_on)
		 ON DUPLICATE KEY UPDATE name = :name, variant = :variant, created_on = :created_on`,
		recipe,
	)
	if err != nil {
		return 0, err
	}

	return recipe.ID, nil

}

func (c *Client) ExistsRecipe(ctx context.Context, recipeID uint64) (uint64, error) {
	if recipeID == 0 {
		return 0, ErrInvalidRecipe
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM recipe WHERE id = ?`,
		recipeID,
	)
	if err := row.Scan(&recipeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecipeNotFound
		}
		return 0, err
	}

	return recipeID, nil
}

func (c *Client) ExistsRecipeName(ctx context.Context, name, variant string) (uint64, error) {
	if name == "" && variant == "" {
		return 0, nil
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM recipe WHERE name = ? AND variant = ?`,
		name,
		variant,
	)
	var recipeID uint64
	if err := row.Scan(&recipeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecipeNotFound
		}
		return 0, err
	}
	return recipeID, nil
}
