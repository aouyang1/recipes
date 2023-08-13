package store

import (
	"context"
	"database/sql"
	"errors"
	"recipes/store/models"
	"strings"
)

var (
	ErrInvalidIngredient        = errors.New("empty ingredient name")
	ErrDuplicateIngredient      = errors.New("cannot insert duplicate ingredient")
	ErrIngredientNotFound       = errors.New("ingredient not found")
	ErrIngredientInUseByRecipes = errors.New("ingredient is in use by recipes")
)

func (c *Client) InsertIngredient(ctx context.Context, ingredient *models.Ingredient) error {
	_, err := c.ExistsIngredient(ctx, ingredient.Name)
	if err == nil {
		return ErrDuplicateIngredient
	}
	if !errors.Is(err, ErrIngredientNotFound) {
		return err
	}
	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO ingredient (id, name)
			  VALUES (:id, :name)`,
		ingredient,
	)
	return err
}

func (c *Client) ExistsIngredient(ctx context.Context, name string) (uint64, error) {
	if name == "" {
		return 0, ErrInvalidIngredient
	}

	var ingredientID uint64
	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM ingredient WHERE name = ?`,
		strings.ToLower(name),
	)
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
		`SELECT id, name FROM ingredient`,
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

func (c *Client) DeleteIngredient(ctx context.Context, name string) error {
	if name == "" {
		return ErrInvalidIngredient
	}

	recipes, err := c.GetIngredientRecipes(ctx, []string{strings.ToLower(name)})
	if err != nil {
		return err
	}

	if len(recipes) > 0 {
		return ErrIngredientInUseByRecipes
	}

	_, err = c.conn.ExecContext(
		ctx,
		`DELETE FROM ingredient WHERE name = ?`,
		strings.ToLower(name),
	)
	return err
}
