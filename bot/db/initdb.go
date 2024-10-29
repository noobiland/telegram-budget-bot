package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"telegram-budget-bot/bot/util"

	_ "github.com/mattn/go-sqlite3"
)

func Init() {
	// Connect to SQLite database (or create if not exists)
	db, err := sql.Open("sqlite3", "./output/expenses.db")
	if err != nil {
		util.Logger.Error("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Check and create table if necessary
	tableName := "expenses"
	err = createTableIfNotExists(db, tableName)
	if err != nil {
		util.Logger.Error("Can't create table", "error", err)
	}
}

// createTableIfNotExists checks if a table exists, and creates it if it does not
func createTableIfNotExists(db *sql.DB, tableName string) error {
	// Query to check if the table exists
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	// If table does not exist, create it
	if err == sql.ErrNoRows {
		createTableQuery := fmt.Sprintf(`
			CREATE TABLE %s (
		"timestamp"	INTEGER NOT NULL,
		"user"	TEXT NOT NULL,
		"amount"	INTEGER NOT NULL,
		"category"	TEXT NOT NULL,
		"payment"	TEXT NOT NULL
	);`, tableName)

		_, err := db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
		slog.Debug(fmt.Sprintf("Table '%s' created successfully.\n", tableName))
	} else {
		slog.Debug(fmt.Sprintf("Table '%s' already exists.\n", tableName))
	}

	return nil
}