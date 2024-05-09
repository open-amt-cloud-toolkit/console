// Package postgres implements postgres connection.
package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	// "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second

	UniqueViolation = "23505"
)

// SQL -.
type SQL struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder    squirrel.StatementBuilderType
	Pool       *sql.DB
	IsEmbedded bool
}

// New -.
func New(url string, opts ...Option) (*SQL, error) {
	pg := &SQL{
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
	if strings.HasPrefix(url, "postgres://") {
		pg.Pool, err = sql.Open("pgx", url)
		if err != nil {
			return nil, err
		}
	} else {
		pg.IsEmbedded = true
		pg.Pool, err = sql.Open("sqlite3", "./console.db")
		if err != nil {
			return nil, err
		}
	}

	return pg, nil
}

// Close -.
func (p *SQL) Close() {
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
