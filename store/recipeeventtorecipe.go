package store

import (
	"context"
	"errors"

	"recipes/store/models"
)

var (
	ErrInvalidRecipeEventID        = errors.New("empty recipe event id")
	ErrRecipeEventToRecipeNotFound = errors.New("recipe event to recipe not found")
)

func (c *Client) UpsertRecipeEventToRecipe(ctx context.Context, recipeEventID string, recipeID uint64) error {
	if recipeEventID == "" {
		return ErrInvalidRecipeEventID
	}

	if recipeID == 0 {
		return ErrInvalidRecipe
	}

	_, err := c.ExistsRecipeEvent(ctx, recipeEventID)
	if err != nil {
		return err
	}

	if _, err := c.ExistsRecipe(ctx, recipeID); err != nil {
		return err
	}

	re2r := models.RecipeEventToRecipe{
		RecipeEventID: recipeEventID,
		RecipeID:      recipeID,
	}

	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_event_to_recipe (recipe_event_id, recipe_id)
			  VALUES (:recipe_event_id, :recipe_id)
		 ON DUPLICATE KEY UPDATE recipe_id = :recipe_event_id, recipe_id = :recipe_id`,
		re2r,
	)
	return err
}

func (c *Client) GetRecipeEventRecipes(ctx context.Context, recipeEventID string) ([]*models.Recipe, error) {
	if recipeEventID == "" {
		return nil, ErrInvalidRecipeEventID
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name, variant, created_on
		   FROM recipe
		   JOIN (SELECT recipe_id
		           FROM recipe_event_to_recipe
		          WHERE recipe_event_id = ?) as re2r
			 ON recipe.id = re2r.recipe_id`,
		recipeEventID,
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

func (c *Client) GetRecipeRecipeEvents(ctx context.Context, name, variant string) ([]*models.RecipeEvent, error) {
	if name == "" {
		return nil, ErrInvalidRecipe
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, title, schedule_date, description 
		   FROM recipe_event
		   JOIN (SELECT recipe_event_id
		           FROM recipe_event_to_recipe
				  WHERE recipe_id = (SELECT id FROM recipe WHERE name = ? AND variant = ?)) as re2r
			 ON recipe_event.id = re2r.recipe_event_id
	   ORDER BY schedule_date DESC`,
		name,
		variant,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.RecipeEvent
	for rows.Next() {
		var id string
		var title string
		var scheduleDate int64
		var description string
		if err := rows.Scan(&id, &title, &scheduleDate, &description); err != nil {
			return nil, err
		}
		res = append(res, &models.RecipeEvent{
			ID:           id,
			Title:        title,
			ScheduleDate: scheduleDate,
			Description:  description,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

type DeleteRecipeEventToRecipeParams struct {
	RecipeEventID string `db:"recipe_event_id"`
	RecipeID      uint64 `db:"recipe_id"`
}

func (c *Client) DeleteRecipeEventToRecipe(ctx context.Context, params *DeleteRecipeEventToRecipeParams) error {
	if params == nil {
		return nil
	}

	if params.RecipeEventID == "" {
		return ErrInvalidRecipeEventID
	}

	if params.RecipeID == 0 {
		return ErrInvalidRecipe
	}

	result, err := c.conn.NamedExecContext(
		ctx,
		`DELETE FROM recipe_event_to_recipe
		       WHERE recipe_event_id = :recipe_event_id
			     AND recipe_id = :recipe_id`,
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
		return ErrRecipeEventToRecipeNotFound
	}
	return nil
}
