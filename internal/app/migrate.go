//go:build migrate

package app

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	dbdbdb "github.com/golang-migrate/migrate/v4/database/sqlite"
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {

	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}
	if strings.HasPrefix(databaseURL, "postgres://") {
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
			log.Fatalf("Migrate: postgres connect error: %s", err)
		}

		err = m.Up()
		defer m.Close()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migrate: up error: %s", err)
		}

		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrate: no change")
			return
		}

		log.Printf("Migrate: up success")
	} else {

		log.Printf("DB path : %s\n", filepath.Join("./", "console.db"))

		db, err := sql.Open("sqlite3", filepath.Join("./", "console.db"))
		if err != nil {
			return
		}
		defer func() {
			if err := db.Close(); err != nil {
				return
			}
		}()
		driver, err := dbdbdb.WithInstance(db, &dbdbdb.Config{})
		if err != nil {
			log.Fatal(err)
		}
		m, err := migrate.NewWithDatabaseInstance("file://migrations", "ql", driver)
		if err != nil {
			log.Fatal(err)
		}
		err = m.Up()
		if err != nil {
			log.Fatal(err)
		}
		defer m.Close()
	}

}
