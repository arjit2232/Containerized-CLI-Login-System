package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"osto-login-cli/internal/appctx"
	"osto-login-cli/internal/cli"
	"osto-login-cli/internal/config"
	"osto-login-cli/internal/db"
	"osto-login-cli/internal/repository"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	pool, err := db.Connect(cfg.DBURL)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer pool.Close()

	if err := runMigrations(pool); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	appCtx := &appctx.AppContext{
		UserRepo:    repository.NewUserRepository(pool),
		SessionRepo: repository.NewSessionRepository(pool),
		Config:      cfg,
	}

	return cli.Run(appCtx)
}

// runMigrations applies any pending golang-migrate migrations found in
// ./migrations against the connected MySQL database.
func runMigrations(pool *sql.DB) error {
	driver, err := mysqlmigrate.WithInstance(pool, &mysqlmigrate.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
