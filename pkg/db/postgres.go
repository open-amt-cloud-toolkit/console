// Package postgres implements postgres connection.
package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	_ "github.com/mattn/go-sqlite3"    // sqlite3 driver
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
	db := &SQL{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(db)
	}

	db.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	var err error

	if strings.HasPrefix(url, "postgres://") {
		err = setupHostedDB(db, url)
		if err != nil {
			return nil, err
		}
	} else {
		err = setupEmbeddedDB(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func setupEmbeddedDB(db *SQL) error {
	db.IsEmbedded = true

	dirname, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	db.Pool, err = sql.Open("sqlite3", filepath.Join(dirname, "device-management-toolkit", "console.db"))
	if err != nil {
		return err
	}

	return nil
}

func setupHostedDB(db *SQL, url string) error {
	var err error

	db.Pool, err = sql.Open("pgx", url)
	if err != nil {
		return err
	}

	return nil
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
