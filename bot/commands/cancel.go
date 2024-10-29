package commands

import (
	"context"
	"fmt"
	"log/slog"
	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/mode"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Discard everything
func CancelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}
	slog.Info(fmt.Sprintf("ChatId: %d discarded the data", update.Message.Chat.ID))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Everything was discarded",
	})
	mode.Storage.Lock()
	mode.Storage.M[update.Message.Chat.ID] = mode.InputAmount
	mode.Storage.Unlock()
}
