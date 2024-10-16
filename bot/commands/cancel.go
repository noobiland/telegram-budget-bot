package commands

import (
	"context"
	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/shared"

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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Everything was discarded",
	})
	shared.Mode.Lock()
	shared.Mode.M[update.Message.Chat.ID] = ""
	shared.Mode.Unlock()
}
