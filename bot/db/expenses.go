package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Insert(user string, amount int, category string, payment string) {
	// Open the SQLite database, it will create the file if it doesn't exist
	db, err := sql.Open("sqlite3", "./output/expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	unixTimestamp := time.Now().Unix()

	// Insert a record
	_, err = db.Exec("INSERT INTO expenses (timestamp, user, amount, category, payment) VALUES (?,?,?,?,?)", unixTimestamp, user, amount, category, payment)
	if err != nil {
		log.Fatalf("Failed to insert record: %v", err)
	}

	// Query records
	rows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp, amount int
		var user, category, payment string
		err = rows.Scan(&timestamp, &user, &amount, &category, &payment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("timestamp: %d, user: %s, amount: %d, category: %s, payment: %s\n",
			timestamp, user, amount, category, payment)
	}
}
