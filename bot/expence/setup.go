package expence

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/db"
	"telegram-budget-bot/bot/mode"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var userStatus = struct {
	sync.RWMutex
	m map[int64]transactionStatus
}{m: make(map[int64]transactionStatus)}

func SetExpence(userId int64, sum int) {
	userStatus.Lock()
	userStatus.m[userId] = transactionStatus{
		sum:      int(sum),
		category: "",
		payment:  "",
	}
	userStatus.Unlock()
}

// Category
func ChooseCategoryMsg(ctx context.Context, b *bot.Bot, chatId int64) {
	_, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        "Select Spending Category:",
		ReplyMarkup: categoryReplyKeyboard,
	})
}

func ChoosenCategory(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	messageText := update.Message.Text

	_, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "You selected: " + messageText,
	})
	mode.Storage.Lock()
	mode.Storage.M[chatId] = mode.WaitingForPayment
	mode.Storage.Unlock()

	userStatus.Lock()
	ts := userStatus.m[chatId]

	slog.Debug(fmt.Sprintf("ts before categorySet: %+v", ts))
	ts.setCategory(messageText)
	userStatus.m[chatId] = ts
	slog.Debug(fmt.Sprintf("ts after categorySet: %+v", ts))
	userStatus.Unlock()

	ChoosePaymentMsg(ctx, b, chatId)
}

// Payment
func ChoosePaymentMsg(ctx context.Context, b *bot.Bot, chatId int64) {

	var _, ok = auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        "Select Payment Method:",
		ReplyMarkup: paymentReplyKeyboard,
	})
}

func ChoosenPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	messageText := update.Message.Text

	_, ok := auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "You selected: " + messageText,
	})

	mode.Storage.Lock()
	mode.Storage.M[chatId] = mode.WaitingForConfirmation
	mode.Storage.Unlock()

	userStatus.Lock()
	ts := userStatus.m[chatId]
	slog.Debug(fmt.Sprintf("ChatId: %d", chatId))
	slog.Debug(fmt.Sprintf("ts before setPayment: %+v", ts))
	ts.setPayment(messageText)
	userStatus.m[chatId] = ts
	slog.Debug(fmt.Sprintf("ts after setPayment: %+v", ts))
	userStatus.Unlock()

	ConfirmationMsg(ctx, b, chatId)
}

// default button to remove the payment one

func defaultButton(ctx context.Context, b *bot.Bot, update *models.Update) {
	var chatId = update.Message.Chat.ID

	var _, ok = auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Connection successful!",
	})
}

// Confirmation
func ConfirmationMsg(ctx context.Context, b *bot.Bot, chatId int64) {

	var _, ok = auth.GetUserName(chatId)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, chatId)
		return
	}

	response, exists := userStatus.m[chatId]
	slog.Info(fmt.Sprintf("ChatId: %d confirmation response: %+v", chatId, response))
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Something went wrong with data!",
		})
	}
	var ts = userStatus.m[chatId]
	userName, _ := auth.GetUserName(chatId)
	db.Insert(userName, ts.sum, ts.category, ts.payment)

	var confirmationMsg = fmt.Sprintf("Amount:\t\t\t¥%d\nCategory:\t\t%s\nPayment:\t\t%s", response.sum, response.category, response.payment)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        confirmationMsg,
		ReplyMarkup: defaultReplyKeyboard,
	})

	userStatus.Lock()
	slog.Debug(fmt.Sprintf("ChatId: %d", chatId))
	slog.Debug(fmt.Sprintf("ts before reset: %+v", ts))
	ts.reset()
	userStatus.m[chatId] = ts
	slog.Debug(fmt.Sprintf("ts after reset: %+v", ts))
	userStatus.Unlock()

	mode.Storage.Lock()
	mode.Storage.M[chatId] = mode.InputAmount
	mode.Storage.Unlock()
}

// TODO: add confirmation step
func Confirmed(ctx context.Context, b *bot.Bot, chatId int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Expence has confirmed",
	})
}
