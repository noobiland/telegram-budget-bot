package commands

import (
	"context"
	"fmt"

	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Turns on money mode
func PrevStatHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	var name, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	var statMsg = db.PrevMonthStat(name)
	fmt.Print(statMsg)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   statMsg,
	})
}
