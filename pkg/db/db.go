package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(256) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT '');

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	// check if db file exists
	var install bool
	_, err := os.Stat(dbFile)
	if errors.Is(err, os.ErrNotExist) {
		install = true
	} else if err != nil {
		return fmt.Errorf("db existence check: %w", err)
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return fmt.Errorf("db connection: %w", err)
	}

	if install {
		_, err = conn.Exec(schema)
		if err != nil {
			conn.Close()
			return fmt.Errorf("db creation: %w", err)
		}
	}

	db = conn
	return nil
}

func Close() {
	db.Close()
}
