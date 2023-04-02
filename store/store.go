package store

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoConfig = errors.New("no mysql store config provided")
)

type Client struct {
	conn *sqlx.DB
}

func NewClient(cfg *mysql.Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrNoConfig
	}
	conn, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}
