package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlite3(file string) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	if db, err = sql.Open("sqlite3", file); err != nil {
		return nil, err
	}
	if err = enableWriteAheadLogging(db); err != nil {
		return nil, err
	}
	return db, nil
}

func enableWriteAheadLogging(db *sql.DB) error {
	if _, err := db.Exec(enableWALSql); err != nil {
		return err
	}
	return nil
}
