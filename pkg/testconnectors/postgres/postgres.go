package testconnectors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewTestPgConnector(ctx context.Context, dsn, migrationDir string) (*pgxpool.Pool, error) {
	parseConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %v", err)
	}

	driver, err := postgres.WithInstance(stdlib.OpenDB(*parseConfig.ConnConfig), &postgres.Config{
		StatementTimeout: time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %v", err)
	}

	migrationsPath := fmt.Sprintf("file://%s", migrationDir)
	mig, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance: %v", err)
	}

	if err = mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("m.Up: %v", err)
	}

	return conn, nil
}
