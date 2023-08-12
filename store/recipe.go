package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"recipes/store/models"
)

var (
	ErrInvalidRecipe   = errors.New("empty recipe name")
	ErrDuplicateRecipe = errors.New("cannot insert duplicate recipe")
	ErrRecipeNotFound  = errors.New("recipe not found")
)

func (c *Client) InsertRecipe(ctx context.Context, recipe *models.Recipe) error {
	_, err := c.ExistsRecipe(ctx, recipe.Name, recipe.Variant)
	if err == nil {
		return ErrDuplicateRecipe
	}
	if !errors.Is(err, ErrRecipeNotFound) {
		return err
	}
	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe (id, name, variant, created_on)
			  VALUES (:id, :name, :variant, :created_on)`,
		recipe,
	)
	return err
}

func (c *Client) ExistsRecipe(ctx context.Context, name, variant string) (uint64, error) {
	if name == "" {
		return 0, ErrInvalidRecipe
	}

	var recipeID uint64
	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM recipe WHERE name = ? AND variant = ?`,
		strings.ToLower(name),
		strings.ToLower(variant),
	)
	if err := row.Scan(&recipeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecipeNotFound
		}
		return 0, err
	}

	return recipeID, nil
}
