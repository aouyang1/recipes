package store

import (
	"context"
	"database/sql"
	"errors"

	"recipes/store/models"
)

var (
	ErrInvalidTag        = errors.New("empty tag name")
	ErrDuplicateTag      = errors.New("cannot insert duplicate tag")
	ErrTagNotFound       = errors.New("tag not found")
	ErrTagInUseByRecipes = errors.New("tag is in use by recipes")
)

func (c *Client) UpsertTag(ctx context.Context, tag *models.Tag) (uint64, error) {
	tid, err := c.ExistsTagName(ctx, tag.Name)
	if err != nil {
		if !errors.Is(err, ErrTagNotFound) {
			return 0, err
		}
	}
	// If the tag name exists, return early. we don't want too update the id
	// of the tag.
	if tid != 0 && tag.ID != tid {
		return tid, ErrDuplicateTag
	}

	if tag.ID == 0 {
		result, err := c.conn.NamedExecContext(
			ctx,
			`INSERT INTO tag (name)
			  VALUES (:name)`,
			tag,
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
		`INSERT INTO tag (id, name)
			  VALUES (:id, :name)
		 ON DUPLICATE KEY UPDATE name = :name`,
		tag,
	)
	if err != nil {
		return 0, err
	}

	return tag.ID, nil
}

func (c *Client) ExistsTag(ctx context.Context, tagID uint64) (uint64, error) {
	if tagID == 0 {
		return 0, ErrInvalidTag
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM tag WHERE id = ?`,
		tagID,
	)
	if err := row.Scan(&tagID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTagNotFound
		}
		return 0, err
	}
	return tagID, nil
}

func (c *Client) ExistsTagName(ctx context.Context, tag string) (uint64, error) {
	if tag == "" {
		return 0, nil
	}

	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM tag WHERE name = ?`,
		tag,
	)
	var tagID uint64
	if err := row.Scan(&tagID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTagNotFound
		}
		return 0, err
	}
	return tagID, nil
}

func (c *Client) GetTags(ctx context.Context) ([]*models.Tag, error) {
	rows, err := c.conn.QueryxContext(
		ctx,
		`SELECT id, name, count(recipe_to_tag.recipe_id) as cnt
		   FROM tag
	  LEFT JOIN recipe_to_tag
		     ON tag.id = recipe_to_tag.tag_id
	   GROUP BY tag.id
	   ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.StructScan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (c *Client) DeleteTag(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrInvalidTag
	}

	recipes, err := c.GetTagRecipes(ctx, id)
	if err != nil {
		return err
	}

	if len(recipes) > 0 {
		return ErrTagInUseByRecipes
	}

	_, err = c.conn.ExecContext(
		ctx,
		`DELETE FROM tag WHERE id = ?`,
		id,
	)
	return err
}
