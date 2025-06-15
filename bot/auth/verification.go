package auth

import (
	"context"
	"telegram-budget-bot/bot/util"

	"github.com/go-telegram/bot"
)

func GetUserName(chatId int64) (string, bool) {
	name, ok := Ids[int(chatId)]
	return name, ok
}

func SendMessageToUnregisteredUser(ctx context.Context, b util.Messenger, chatId int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "This is a private bot. You are not authorized to use it",
	})
}
