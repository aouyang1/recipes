package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"recipes/store/models"
)

var (
	ErrInvalidTag        = errors.New("empty tag name")
	ErrDuplicateTag      = errors.New("cannot insert duplicate tag")
	ErrTagNotFound       = errors.New("tag not found")
	ErrTagInUseByRecipes = errors.New("tag is in use by recipes")
)

func (c *Client) InsertTag(ctx context.Context, tag *models.Tag) error {
	_, err := c.ExistsTag(ctx, tag.Name)
	if err == nil {
		return ErrDuplicateTag
	}
	if !errors.Is(err, ErrTagNotFound) {
		return err
	}

	tag.Name = strings.ToLower(tag.Name)

	_, err = c.conn.NamedExecContext(
		ctx,
		`INSERT INTO tag (id, name)
			  VALUES (:id, :name)`,
		tag,
	)
	return err
}

func (c *Client) ExistsTag(ctx context.Context, name string) (uint64, error) {
	if name == "" {
		return 0, ErrInvalidTag
	}

	var tagID uint64
	row := c.conn.QueryRowxContext(
		ctx,
		`SELECT id FROM tag WHERE name = ?`,
		strings.ToLower(name),
	)
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
		`SELECT id, name FROM tag`,
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

func (c *Client) DeleteTag(ctx context.Context, name string) error {
	if name == "" {
		return ErrInvalidTag
	}

	recipes, err := c.GetTagRecipes(ctx, []string{strings.ToLower(name)})
	if err != nil {
		return err
	}

	if len(recipes) > 0 {
		return ErrTagInUseByRecipes
	}

	_, err = c.conn.ExecContext(
		ctx,
		`DELETE FROM tag WHERE name = ?`,
		strings.ToLower(name),
	)
	return err
}
