package db

import (
	"database/sql"
	"telegram-budget-bot/bot/util"
)

func GetUsers() map[int]string {
	db, err := sql.Open("sqlite3", "./output/users.db")
	if err != nil {
		util.Logger.Error("cant open db: ", "error", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		util.Logger.Error("cant query users table: ", "error", err)
	}
	defer rows.Close()
	var idsMap = make(map[int]string)
	for rows.Next() {
		var id, telegram_chat_id int
		var name string
		err = rows.Scan(&id, &name, &telegram_chat_id)
		if err != nil {
			util.Logger.Error("cant scan for rows in table", "error", err)
		}
		idsMap[telegram_chat_id] = name
	}
	return idsMap
}
