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

	_, err := c.conn.NamedExecContext(
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

func (c *Client) ExistsRecipe(ctx context.Context, name, variant string) (uint64, error) {
	if name == "" {
		return 0, ErrInvalidRecipe
	}

	var recipeID uint64
	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM recipe WHERE name = ? AND variant = ?`,
		name,
		variant,
	)
	if err := row.Scan(&recipeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecipeNotFound
		}
		return 0, err
	}

	return recipeID, nil
}
