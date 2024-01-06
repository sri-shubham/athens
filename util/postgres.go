package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var db *pg.DB

// CustomQueryLogger implements the pg.QueryHook interface
type CustomQueryLogger struct{}

func (c *CustomQueryLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	// Log the query before execution
	query, _ := q.FormattedQuery()
	log.Printf("Executing query: %s", query)
	return ctx, nil
}

func (c *CustomQueryLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	// Log any errors or other relevant information after query execution
	if q.Err != nil {
		log.Printf("Query error: %s", q.Err)
	}
	return nil
}

func ConnectToPostgres(config *Config) error {
	db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Postgres.Host, config.Postgres.Port),
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		Database: config.Postgres.DBName,
	})

	if config.Postgres.LogQueries {
		db.AddQueryHook(&CustomQueryLogger{})
	}

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
