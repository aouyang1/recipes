package store

import (
	"context"
	"errors"
	"recipes/store/models"
)

var (
	ErrNilRecipeEvent = errors.New("nil recipe event")
)

func (c *Client) UpsertRecipeEventContext(ctx context.Context, recipeEvent *models.RecipeEvent) error {
	if recipeEvent == nil {
		return ErrNilRecipeEvent
	}

	exists, err := c.ExistsRecipeEventContext(ctx, recipeEvent)
	if err != nil {
		return err
	}

	// insert new record
	if !exists {
		_, err := c.conn.NamedExecContext(
			ctx,
			`INSERT INTO recipe_event (id, schedule_date, title, description)
		 	  	  VALUES (:id, :schedule_date, :title, :description)`,
			recipeEvent,
		)
		return err
	}

	// update
	_, err = c.conn.NamedExecContext(
		ctx,
		`UPDATE recipe_event
			SET schedule_date = :schedule_date, title = :title, description = :description
		  WHERE id = :id`,
		recipeEvent,
	)
	return err
}

func (c *Client) ExistsRecipeEventContext(ctx context.Context, recipeEvent *models.RecipeEvent) (bool, error) {
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
