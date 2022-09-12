package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pool interface {
	pgxtype.Querier
}

type ConnectionConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	SSL      string
}

// NewPool initializes a new postgres connection pool.
func NewPool(c ConnectionConfig) (*pgxpool.Pool, error) {
	postgresConn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DB, c.SSL,
	)

	pgxConfig, err := pgxpool.ParseConfig(postgresConn)
	if err != nil {
		return nil, fmt.Errorf("error parsing postgres config %w", err)
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("error initializing connection pool %w", err)
	}

	return dbPool, nil
}
