package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var mode = struct {
	sync.RWMutex
	m map[int64]string
}{m: make(map[int64]string)}

func main() {
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

	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, cancelMode)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/money", bot.MatchTypeExact, moneyMode)
	
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	mode.RLock()
	state, exists := mode.m[update.Message.Chat.ID]
	mode.RUnlock()

	if !exists || state == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Use Command",
		})
	} else if state == "money" {
		handlerLogic(ctx, b, update)
	}

}

// Discard everything
func cancelMode(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Everything was discarded",
	})
	mode.Lock()
	mode.m[update.Message.Chat.ID] = ""
	mode.Unlock()
}

// Turns on money mode
func moneyMode(ctx context.Context, b *bot.Bot, update *models.Update) {
	mode.RLock()
	state, exists := mode.m[update.Message.Chat.ID]
	mode.RUnlock()

	if !exists || state == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Enter the sum",
		})
		mode.Lock()
		mode.m[update.Message.Chat.ID] = "money"
		mode.Unlock()
	}
}

// Actual handler
func handlerLogic(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "dummy",
	})
	mode.Lock()
	mode.m[update.Message.Chat.ID] = ""
	mode.Unlock()
}
