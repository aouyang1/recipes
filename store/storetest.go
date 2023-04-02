package store

import (
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type TestClient struct {
	conn *sqlx.DB
}

func NewTestClient() (*TestClient, error) {
	client, err := NewClient(TestConfig())
	if err != nil {
		return nil, err
	}
	return &TestClient{
		conn: client.conn,
	}, nil
}

func TestConfig() *mysql.Config {
	return &mysql.Config{
		User:   os.Getenv("USER_MYSQL_TEST_USERNAME"),
		Passwd: os.Getenv("USER_MYSQL_TEST_PASSWORD"),
		DBName: os.Getenv("USER_MYSQL_TEST_DB"),
		Net:    "tcp",
		Addr:   "localhost",
	}
}

func (t *TestClient) Cleanup() error {
	tx := t.conn.MustBegin()
	tx.MustExec("TRUNCATE TABLE recipe_event")
	tx.MustExec("TRUNCATE TABLE recipe_event_to_recipe")
	tx.MustExec("TRUNCATE TABLE recipe")
	tx.MustExec("TRUNCATE TABLE recipe_tag")
	tx.MustExec("TRUNCATE TABLE ingredient")
	tx.MustExec("TRUNCATE TABLE recipe_ingredient")
	return tx.Commit()
}
