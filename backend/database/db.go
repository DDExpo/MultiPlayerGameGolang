package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type UserData struct {
	ID             string
	Username       string
	TimeRegistered int64
}

func GetDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "./store.db")
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func DbDelete(DB *sqlx.DB, ids []string, table string) error {

	query, args, err := sqlx.In("DELETE FROM "+table+" WHERE id IN (?)", ids)
	if err != nil {
		return err
	}
	query = DB.Rebind(query)
	_, err = DB.Exec(query, args...)
	return err
}

func DBGetUser(DB *sqlx.DB, sessionID string) (string, int64, error) {
	row := DB.QueryRow("SELECT id, username, timeRegistered FROM users WHERE id = ?", sessionID)

	var id string
	var username string
	var timeRegistered int64
	err := row.Scan(&id, &username, &timeRegistered)
	if err != nil {
		return "", 0, fmt.Errorf("No user found with that id: %w", err)
	}
	return username, timeRegistered, nil
}

func DBSaveUser(DB *sqlx.DB, user UserData) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO users (id, username, timeRegistered)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			username = excluded.username,
			timeRegistered = excluded.timeRegistered
	`, user.ID, user.Username, user.TimeRegistered)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error saving user: %v", err)
	}

	return tx.Commit()
}

func RunFirstTimeShemas(db *sqlx.DB) error {
	schemaUsers := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT,
		timeRegistered BIGINT
	);`
	_, err := db.Exec(schemaUsers)
	if err != nil {
		return fmt.Errorf("error on executing schemaUsers: %v", err)
	}
	return nil
}
