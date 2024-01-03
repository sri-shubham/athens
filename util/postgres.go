package util

import (
	"fmt"
	"io"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func GetDb() *pg.DB {
	return db
}

func InitPostgresDB() error {
	zap.L().Info("Starting database initialization (postgres)")
	defer zap.L().Info("Completed database initialization (postgres)")

	file, err := os.Open("sql/init.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	queries, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(queries))
	if err != nil {
		return err
	}

	return nil
}
