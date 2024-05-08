// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver i think
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second

	UniqueViolation = "23505"
)

// DB -.
type DB struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *sql.DB
}

// New -.
func New(url string, opts ...Option) (*DB, error) {
	pg := &DB{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	var err error

	pg.Pool, err = sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

// Close -.
func (p *DB) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func CheckNotUnique(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == UniqueViolation {
			return true
		}
	}

	return false
}
