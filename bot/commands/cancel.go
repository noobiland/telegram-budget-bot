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
	chatId := update.Message.Chat.ID
	_, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}
	slog.Info(fmt.Sprintf("ChatId: %d discarded the data", chatId))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Everything was discarded",
	})
	mode.Storage.Lock()
	mode.Storage.M[chatId] = mode.InputAmount
	mode.Storage.Unlock()
}
