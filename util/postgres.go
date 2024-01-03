package util

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

var db *pg.DB

func ConnectToPostgres(config *Config) error {
	db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Postgres.Host, config.Postgres.Port),
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		Database: config.Postgres.DBName,
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return errors.Wrap(err, "Error connecting to PostgreSQL")
	}

	return nil
}
