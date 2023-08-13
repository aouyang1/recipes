package store

import (
	"context"
	"errors"
	"strings"

	"recipes/store/models"
)

var (
	ErrRecipeToTagNotFound = errors.New("recipe to tag not found")
)

func (c *Client) UpsertRecipeToTag(ctx context.Context, recipeName, recipeVariant, tag string) error {
	if recipeName == "" {
		return ErrInvalidRecipe
	}
	if tag == "" {
		return ErrInvalidTag
	}

	recipeID, err := c.ExistsRecipe(ctx, recipeName, recipeVariant)
	if err != nil {
		return err
	}

	tagID, err := c.ExistsTag(ctx, tag)
	if err != nil {
		return err
	}

	r2t := models.RecipeToTag{
		RecipeID: recipeID,
		TagID:    tagID,
	}

	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO recipe_to_tag (recipe_id, tag_id)
			  VALUES (:recipe_id, :tag_id)
		 ON DUPLICATE KEY UPDATE recipe_id = :recipe_id, tag_id = :tag_id`,
		r2t,
	)
	return err
}

func (c *Client) GetRecipeTags(ctx context.Context, name, variant string) ([]*models.Tag, error) {
	if name == "" {
		return nil, ErrInvalidRecipe
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name
		   FROM tag
		   JOIN (SELECT tag_id
		           FROM recipe_to_tag
		          WHERE recipe_id = (SELECT id FROM recipe WHERE name = ? AND variant = ?) as r2t
			 ON tag.id = r2t.tag_id`,
		strings.ToLower(name),
		strings.ToLower(variant),
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

func (c *Client) GetTagRecipes(ctx context.Context, tags []string) ([]*models.Recipe, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	for i, tag := range tags {
		tags[i] = strings.ToLower(tag)
	}

	rows, err := c.conn.QueryContext(
		ctx,
		`SELECT id, name, variant, created_on
		   FROM recipe 
		   JOIN (SELECT recipe_id
		           FROM recipe_to_tag
		           JOIN (SELECT id FROM tag WHERE name IN (?))
				     ON recipe_to_tag.tag_id = tag.id) as r2t
			 ON recipe.id = r2t.recipe_id`,
		tags,
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

type DeleteRecipeToTagParams struct {
	RecipeName    string `db:"recipe_name"`
	RecipeVariant string `db:"recipe_variant"`
	Tag           string `db:"tag"`
}

func (c *Client) DeleteRecipeToTag(ctx context.Context, params DeleteRecipeToTagParams) error {
	if params.RecipeName == "" {
		return ErrInvalidRecipe
	}
	if params.Tag == "" {
		return ErrInvalidTag
	}

	result, err := c.conn.NamedExecContext(
		ctx,
		`DELETE FROM recipe_to_tag
		       WHERE recipe_id = (SELECT id FROM recipe WHERE name = :recipe_name AND variant = :recipe_variant)
			     AND tag_id = (SELECT id FROM tag WHERE name = :tag)`,
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
