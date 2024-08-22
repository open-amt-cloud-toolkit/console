package sqldb

import "database/sql"

// SQL is an interface that abstracts the SQL database operations.
type SQL interface {
	// Define the methods that `db.SQL` should have.
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SQLWrapper wraps *sql.DB to implement the db.SQL interface.
type SQLWrapper struct {
	DB *sql.DB
}

func (sw *SQLWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	return sw.DB.QueryRow(query, args...)
}

func (sw *SQLWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return sw.DB.Query(query, args...)
}

func (sw *SQLWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return sw.DB.Exec(query, args...)
}
