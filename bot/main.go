package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/commands"
	"telegram-budget-bot/bot/db"
	"telegram-budget-bot/bot/expence"
	"telegram-budget-bot/bot/shared"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	db.Init()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	token, err := os.ReadFile("resources/token")
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(string(token), opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, commands.HelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, commands.CancelHandler)
	expence.InitDefaultKeyboard(b)
	expence.InitCategoryKeyboard(b)
	expence.InitPaymentKeyboard(b)

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Got a message from chat Id: ", update.Message.Chat.ID)

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	shared.Mode.RLock()
	state, exists := shared.Mode.M[update.Message.Chat.ID]
	shared.Mode.RUnlock()

	if !exists || state == "" {
		sum, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Invalid input, %s is not a number", update.Message.Text),
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Valid input, %d is number!!!", sum),
			})
			expence.SetExpence(update.Message.Chat.ID, sum)
			expence.ChooseCategoryMsg(ctx, b, update)
			shared.Mode.Lock()
			shared.Mode.M[update.Message.Chat.ID] = "waiting_for_category"
			shared.Mode.Unlock()
		}
	} else if state == "waiting_for_category" {
		expence.ChooseCategoryMsg(ctx, b, update)
	} else if state == "waiting_for_payment" {
		expence.ChoosePaymentMsg(ctx, b, update)
	} else if state == "waiting_for_confirmation" {
		expence.ConfirmationMsg(ctx, b, update)
	}
}
