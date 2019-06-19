package database

import (
	"database/sql"
)

func ExecTransaction(db *sql.DB, query string, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
