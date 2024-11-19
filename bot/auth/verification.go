package auth

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetUserName(chatId int64) (string, bool) {
	name, ok := Ids[int(chatId)]
	return name, ok
}

func SendMessageToUnregisteredUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This is a private bot. You are authorized to use it",
	})
}
