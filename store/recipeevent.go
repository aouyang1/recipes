package store

import (
	"context"
	"errors"
	"recipes/store/models"
)

var (
	ErrNilRecipeEvent = errors.New("nil recipe event")
)

func (c *Client) UpsertRecipeEvent(ctx context.Context, recipeEvent *models.RecipeEvent) error {
	if recipeEvent == nil {
		return ErrNilRecipeEvent
	}

	_, err := c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_event (id, schedule_date, title, description)
		 	  	  VALUES (:id, :schedule_date, :title, :description)
		 ON DUPLICATE KEY UPDATE schedule_date = :schedule_date, title = :title, description = :description`,
		recipeEvent,
	)
	return err
}

func (c *Client) ExistsRecipeEvent(ctx context.Context, recipeEvent *models.RecipeEvent) (bool, error) {
	if recipeEvent == nil {
		return false, ErrNilRecipeEvent
	}

	var cnt int
	err := c.conn.QueryRowxContext(
		ctx,
		`SELECT COUNT(id) FROM recipe_event WHERE id = ?`,
		recipeEvent.ID,
	).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
