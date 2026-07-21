package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Connect opens a MySQL connection pool and retries until the DB is
// reachable (useful when the DB container is still starting up) or the
// retry budget is exhausted.
func Connect(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(5 * time.Minute)

	const maxAttempts = 10
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = pool.Ping()
		if lastErr == nil {
			return pool, nil
		}
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	return nil, fmt.Errorf("could not connect to db after %d attempts: %w", maxAttempts, lastErr)
}
