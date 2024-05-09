package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	dbdbdb "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for file source
	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultAttempts     = 20
	_defaultTimeout      = time.Second
	_directoryPermission = 0o755
)

var errMigrate = errors.New("migrate error")

func MigrationError(op string) error {
	return fmt.Errorf("%w: %s", errMigrate, op)
}

func Init() error {
	databaseURL, ok := os.LookupEnv("DB_URL")
	if !ok || databaseURL == "" {
		log.Printf("migrate: environment variable not declared: DB_URL -- using embedded database")
	}

	if strings.HasPrefix(databaseURL, "postgres://") {
		err := setupHostedDB(databaseURL)
		if err != nil {
			return err
		}
	} else {
		// make sure the directory exists
		err := setupLocalDB()
		if err != nil {
			return err
		}
	}

	return nil
}

func setupLocalDB() error {
	dirname, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	consoleDir := filepath.Join(dirname, "device-management-toolkit")

	if _, err = os.Stat(consoleDir); os.IsNotExist(err) {
		if err1 := os.Mkdir(consoleDir, _directoryPermission); err1 != nil {
			return err1
		}
	}

	log.Printf("DB path : %s\n", filepath.Join(consoleDir, "console.db"))

	db, err := sql.Open("sqlite3", filepath.Join(consoleDir, "console.db"))
	if err != nil {
		return err
	}

	defer func() {
		if err1 := db.Close(); err1 != nil {
			return
		}
	}()

	driver, err := dbdbdb.WithInstance(db, &dbdbdb.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "ql", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	defer m.Close()

	return nil
}

func setupHostedDB(databaseURL string) error {
	databaseURL += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)

		attempts--
	}

	if err != nil {
		return MigrationError(fmt.Sprintf("postgres connect error: %s", err))
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return MigrationError(fmt.Sprintf("up error: %s", err))
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")

		return nil
	}

	log.Printf("Migrate: up success")

	return nil
}
