package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"recipes/store/models"
)

var (
	ErrInvalidTag   = errors.New("empty tag name")
	ErrDuplicateTag = errors.New("cannot insert duplicate tag")
	ErrTagNotFound  = errors.New("tag not found")
)

func (c *Client) InsertTag(ctx context.Context, tag *models.Tag) error {
	_, err := c.ExistsTag(ctx, tag.Name)
	if err == nil {
		return ErrDuplicateTag
	}
	if !errors.Is(err, ErrTagNotFound) {
		return err
	}
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
