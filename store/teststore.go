package store

import (
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type TestClient struct {
	conn *sqlx.DB
}

func NewTestClient() (*Client, error) {
	return NewClient(TestConfig())
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

func Cleanup(conn *sqlx.DB) error {
	tx := conn.MustBegin()
	tx.MustExec("TRUNCATE TABLE recipe_event_to_recipe")
	tx.MustExec("TRUNCATE TABLE recipe_to_tag")
	tx.MustExec("TRUNCATE TABLE recipe_to_ingredient")
	tx.MustExec("TRUNCATE TABLE recipe_event")
	tx.MustExec("TRUNCATE TABLE recipe")
	tx.MustExec("TRUNCATE TABLE tag")
	tx.MustExec("TRUNCATE TABLE ingredient")
	return tx.Commit()
}
