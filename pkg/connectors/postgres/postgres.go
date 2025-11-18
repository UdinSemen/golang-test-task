package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type DatabaseConfig struct {
	Host string `env:"HOST" yaml:"host" validate:"required"`
	Port int    `env:"PORT" yaml:"port" validate:"required"`
	User string `env:"USER" yaml:"user" validate:"required"`
	Pass string `env:"PASS" yaml:"pass" validate:"required"`
	DB   string `env:"DB"   yaml:"db" validate:"required"`
}

type MigrationLogger struct {
	logger  *log.Logger
	verbose bool
}

func NewMigrationLogger(logger *log.Logger, verbose bool) *MigrationLogger {
	return &MigrationLogger{
		logger:  logger,
		verbose: verbose,
	}
}

func (l *MigrationLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Println(msg, "component=migration")
}
func (l *MigrationLogger) Verbose() bool {
	return l.verbose
}

type Connector struct {
	conn *pgxpool.Pool
}

func NewConnector(cfg *DatabaseConfig) (*Connector, error) {
	const op = "NewConnector"
	parseConfig, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s dbname=%s host=%s port=%v password=%s sslmode=disable pool_max_conns=10 pool_max_conn_lifetime=1h30m",
		cfg.User,
		cfg.DB,
		cfg.Host,
		cfg.Port,
		cfg.Pass,
	))
	if err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "pgx.ParseConfig", err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "pgx.ConnectConfig", err)
	}

	driver, err := postgres.WithInstance(stdlib.OpenDB(*parseConfig.ConnConfig), &postgres.Config{
		StatementTimeout: time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "postgres.WithInstance", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "migrate.NewWithDatabaseInstance", err)
	}
	m.Log = NewMigrationLogger(log.Default(), true)
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s.%s: %w", op, "m.Up", err)
	}
	version, _, vErr := m.Version()
	if vErr == nil {
		log.Println("migration completed successfully",
			"version", version,
		)
	}
	return &Connector{conn: conn}, nil

}

func (c *Connector) GetConnector() *pgxpool.Pool {
	return c.conn
}

func (c *Connector) Close() error {
	c.conn.Close()
	return nil
}
