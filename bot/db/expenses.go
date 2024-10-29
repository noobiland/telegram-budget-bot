package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"telegram-budget-bot/bot/util"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Insert(user string, amount int, category string, payment string) {

	// Open the SQLite database, it will create the file if it doesn't exist
	db, err := sql.Open("sqlite3", "./output/expenses.db")
	if err != nil {
		util.Logger.Error("cant open db: ", "error", err)
	}
	defer db.Close()

	unixTimestamp := time.Now().Unix()

	// Insert a record
	_, err = db.Exec("INSERT INTO expenses (timestamp, user, amount, category, payment) VALUES (?,?,?,?,?)", unixTimestamp, user, amount, category, payment)
	if err != nil {
		util.Logger.Error("Failed to insert record: %v", "error", err)
	}
}

func queryData(db *sql.DB){
	// Query records
	rows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		util.Logger.Error("cant query expences tabel: ", "error", err)
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp, amount int
		var user, category, payment string
		err = rows.Scan(&timestamp, &user, &amount, &category, &payment)
		if err != nil {
			util.Logger.Error("cant scan for rows in table", "error", err)
		}
		slog.Info(fmt.Sprintf("timestamp: %d, user: %s, amount: %d, category: %s, payment: %s\n",
		timestamp, user, amount, category, payment))
	}
}