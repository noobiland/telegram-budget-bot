package commands

import (
	"context"
	"log/slog"

	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Stat messages
func PrevStatHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	name, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	statMsg := db.PrevMonthStat(name)
	slog.Info(statMsg)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   statMsg,
	})
}

func CurStatHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	name, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	statMsg := db.CurrMonthStat(name)
	slog.Info(statMsg)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   statMsg,
	})
}
