package store

import (
	"context"
	"errors"

	"recipes/store/models"
)

var (
	ErrRecipeToTagNotFound = errors.New("recipe to tag not found")
)

func (c *Client) UpsertRecipeToTag(ctx context.Context, recipeID uint64, tagID uint64) error {
	if recipeID == 0 {
		return ErrInvalidRecipe
	}
	if tagID == 0 {
		return ErrInvalidTag
	}

	if _, err := c.ExistsRecipe(ctx, recipeID); err != nil {
		return err
	}

	if _, err := c.ExistsTag(ctx, tagID); err != nil {
		return err
	}

	r2t := models.RecipeToTag{
		RecipeID: recipeID,
		TagID:    tagID,
	}

	_, err := c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_to_tag (recipe_id, tag_id)
			  VALUES (:recipe_id, :tag_id)
		 ON DUPLICATE KEY UPDATE recipe_id = :recipe_id, tag_id = :tag_id`,
		r2t,
	)
	return err
}

func (c *Client) GetRecipeTags(ctx context.Context, recipeID uint64) ([]*models.Tag, error) {
	if recipeID == 0 {
		return nil, ErrInvalidRecipe
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name
		   FROM tag
		   JOIN (SELECT tag_id
		           FROM recipe_to_tag
		          WHERE recipe_id = ?) as r2t
			 ON tag.id = r2t.tag_id`,
		recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Tag
	for rows.Next() {
		var id uint64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		res = append(res, &models.Tag{
			ID:   id,
			Name: name,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetTagRecipes(ctx context.Context, tagID uint64) ([]*models.Recipe, error) {
	if tagID == 0 {
		return nil, nil
	}

	rows, err := c.conn.QueryxContext(
		ctx,
		`SELECT id, name, variant, created_on
		   FROM recipe 
		   JOIN (SELECT recipe_id
		           FROM recipe_to_tag
				  WHERE tag_id = ?) as r2t
			 ON recipe.id = r2t.recipe_id`,
		tagID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.StructScan(&recipe); err != nil {
			return nil, err
		}
		res = append(res, &recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

type DeleteRecipeToTagParams struct {
	RecipeID uint64 `db:"recipe_id"`
	TagID    uint64 `db:"tag_id"`
}

func (c *Client) DeleteRecipeToTag(ctx context.Context, params *DeleteRecipeToTagParams) error {
	if params == nil {
		return nil
	}
	if params.RecipeID == 0 {
		return ErrInvalidRecipe
	}
	if params.TagID == 0 {
		return ErrInvalidTag
	}

	result, err := c.conn.NamedExecContext(
		ctx,
		`DELETE FROM recipe_to_tag
		       WHERE recipe_id = :recipe_id
			     AND tag_id = :tag_id`,
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
		return ErrRecipeToTagNotFound
	}
	return nil
}
