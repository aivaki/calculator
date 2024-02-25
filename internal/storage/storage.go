package storage

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	DriverName = `sqlite3`
	DataSourceName = `storage.db`

	expressions = `
	CREATE TABLE IF NOT EXISTS Expressions (
		id INTEGER PRIMARY KEY,
		key TEXT NOT NULL,
		expression TEXT NOT NULL,
		result TEXT,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP,
		error_message TEXT
	);`

	operationsTable = `
	CREATE TABLE IF NOT EXISTS Operations (
		id INTEGER PRIMARY KEY,
		operation_type TEXT NOT NULL,
		execution_time INTEGER NOT NULL
	);`

	operationsValues = `
	INSERT INTO Operations (operation_type, execution_time) VALUES
    ('+', 100),
    ('-', 100),
    ('*', 100),
    ('/', 100);`
)

var (
	query = []string{
		expressions,
		operationsTable,
		operationsValues,
	}
)

func New() (db *sql.DB, err error) {
	if _, err = os.Stat(DataSourceName); err == nil {
		return sql.Open(DriverName, DataSourceName)
	}

	db, err = sql.Open(DriverName, DataSourceName)
	if err != nil {
		err = fmt.Errorf("error creating database: %v", err)
		return
	}
	
	for _, e := range query {
		if 	_, err = db.Exec(e); err != nil {
			err = fmt.Errorf("table error: %v", err)
			return
		}
	}
	return
}
