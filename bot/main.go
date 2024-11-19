package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/commands"
	"telegram-budget-bot/bot/expence"
	"telegram-budget-bot/bot/mode"
	"telegram-budget-bot/bot/util"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	slog.Info("Starting configuration...")
	// TODO: add dbs validation check
	auth.InitUsers()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	slog.Info("Reading token...")
	token, err := os.ReadFile("resources/token")
	if err != nil {
		util.Logger.Error("No Token", "error", err)
	}

	slog.Info("Set up bot")
	b, err := bot.New(string(token), opts...)
	if err != nil {
		util.Logger.Error("Can't create bot instance.", "error", err)
		panic(err)
	}

	slog.Info("Preparing handlers and keyboards...")
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, commands.HelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, commands.CancelHandler)
	expence.InitDefaultKeyboard(b)
	expence.InitCategoryKeyboard(b)
	expence.InitPaymentKeyboard(b)

	slog.Info("Starting context...")
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("Got a message from chat Id: %d", update.Message.Chat.ID))

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	mode.Storage.RLock()
	state, exists := mode.Storage.M[update.Message.Chat.ID]
	mode.Storage.RUnlock()

	if !exists {
		state = mode.InputAmount
	}

	switch state {
	case mode.InputAmount:
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Invalid input, %s is not a number", update.Message.Text),
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("¥%d\nPlease provide the spending category", sum),
			})
			expence.SetExpence(update.Message.Chat.ID, sum)
			expence.ChooseCategoryMsg(ctx, b, update)
			mode.Storage.Lock()
			mode.Storage.M[update.Message.Chat.ID] = mode.WaitingForCategory
			mode.Storage.Unlock()
		}
	case mode.WaitingForCategory:
		expence.ChooseCategoryMsg(ctx, b, update)
	case mode.WaitingForPayment:
		expence.ChoosePaymentMsg(ctx, b, update)
	case mode.WaitingForConfirmation:
		expence.ConfirmationMsg(ctx, b, update)
	}
}
