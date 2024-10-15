package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"
)

// Send any text message to the bot after the bot has been started
var demoReplyKeyboard *reply.ReplyKeyboard

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(inlineModeHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
	}

	token, err := os.ReadFile("resources/token")
	if err != nil {
		log.Fatal(err)
	}


	b, err := bot.New(string(token), opts...)
	if err != nil {
		panic(err)
	}
	// TODO: keep for future web-hook interaction
	// params := &bot.DeleteWebhookParams{
	// 	DropPendingUpdates: true, // Drop pending updates
	// }
	// result, err := b.DeleteWebhook(ctx, params)
	// if err != nil {
	// 	panic(err)
	// }
	// if result {
	// 	fmt.Println("Webhook deleted successfully")
	// }
	initReplyKeyboard(b)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/money", bot.MatchTypeExact, handlerReplyKeyboard)
	b.Start(ctx)
}

// callback for buttons
func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}

// buttons near message
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "button_1"},
				{Text: "Button 2", CallbackData: "button_2"},
			}, {
				{Text: "Button 3", CallbackData: "button_3"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Click by button",
		ReplyMarkup: kb,
	})
}

// inline mode
func inlineModeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.InlineQuery == nil {
		return
	}

	results := []models.InlineQueryResult{
		&models.InlineQueryResultArticle{ID: "1", Title: "Foo 1", InputMessageContent: &models.InputTextMessageContent{MessageText: "hui 1"}},
		&models.InlineQueryResultArticle{ID: "2", Title: "Foo 2", InputMessageContent: &models.InputTextMessageContent{MessageText: "hui 2"}},
		&models.InlineQueryResultArticle{ID: "3", Title: "Foo 3", InputMessageContent: &models.InputTextMessageContent{MessageText: "hui 3"}},
	}

	b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
	})
}

// Custom keyboarda
func initReplyKeyboard(b *bot.Bot) {
	demoReplyKeyboard = reply.New(reply.WithPrefix("reply_keyboard"), reply.IsSelective(), reply.IsOneTimeKeyboard()).
		Button("Button", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Row().
		Button("Cancel", b, bot.MatchTypeExact, onReplyKeyboardSelect)
}

func handlerReplyKeyboard(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select example command from reply keyboard:",
		ReplyMarkup: demoReplyKeyboard,
	})
}

func onReplyKeyboardSelect(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You selected: " + string(update.Message.Text),
	})
}
