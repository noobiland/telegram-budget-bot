package auth

import (
	"telegram-budget-bot/bot/db"
)

var Ids = make(map[int]string)

func InitUsers() {
	Ids = db.GetUsers()
}
