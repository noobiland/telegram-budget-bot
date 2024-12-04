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

	var name, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	var statMsg = db.PrevMonthStat(name)
	slog.Info(statMsg)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   statMsg,
	})
}

func CurStatHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	var name, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	var statMsg = db.CurrMonthStat(name)
	slog.Info(statMsg)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   statMsg,
	})
}