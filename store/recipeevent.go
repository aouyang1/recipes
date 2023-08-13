package store

import (
	"context"
	"database/sql"
	"errors"
	"recipes/store/models"
)

var (
	ErrNilRecipeEvent      = errors.New("nil recipe event")
	ErrRecipeEventNotFound = errors.New("recipe event not found")
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

func (c *Client) ExistsRecipeEvent(ctx context.Context, recipeEventID string) (bool, error) {
	if recipeEventID == "" {
		return false, ErrNilRecipeEvent
	}

	var cnt int
	err := c.conn.QueryRowxContext(
		ctx,
		`SELECT COUNT(id) FROM recipe_event WHERE id = ?`,
		recipeEventID,
	).Scan(&cnt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrRecipeEventNotFound
		}
		return false, err
	}
	return cnt > 0, nil
}

func (c *Client) GetRecipeEvents(ctx context.Context) ([]*models.RecipeEvent, error) {
	rows, err := c.conn.QueryxContext(
		ctx,
		`SELECT id, schedule_date, title, description
		   FROM recipe_event
	   ORDER BY schedule_date DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipeEvents []*models.RecipeEvent
	for rows.Next() {
		var recipeEvent models.RecipeEvent
		if err := rows.StructScan(&recipeEvent); err != nil {
			return nil, err
		}
		recipeEvents = append(recipeEvents, &recipeEvent)
	}
	return recipeEvents, nil
}
