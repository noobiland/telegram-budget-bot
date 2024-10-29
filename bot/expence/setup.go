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
func ChooseCategoryMsg(ctx context.Context, b *bot.Bot, update *models.Update) {
	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select Spending Category:",
		ReplyMarkup: categoryReplyKeyboard,
	})
}

func ChoosenCategory(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You selected: " + string(update.Message.Text),
	})
	mode.Storage.Lock()
	mode.Storage.M[update.Message.Chat.ID] = mode.WaitingForPayment
	mode.Storage.Unlock()

	userStatus.Lock()
	var ts = userStatus.m[update.Message.Chat.ID]

	slog.Debug(fmt.Sprintf("ts before categorySet: %+v", ts))
	ts.setCategory(string(update.Message.Text))
	userStatus.m[update.Message.Chat.ID] = ts
	slog.Debug(fmt.Sprintf("ts after categorySet: %+v", ts))
	userStatus.Unlock()

	ChoosePaymentMsg(ctx, b, update)
}

// Payment
func ChoosePaymentMsg(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select Payment Method:",
		ReplyMarkup: paymentReplyKeyboard,
	})
}

func ChoosenPayment(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You selected: " + string(update.Message.Text),
	})

	mode.Storage.Lock()
	mode.Storage.M[update.Message.Chat.ID] = mode.WaitingForConfirmation
	mode.Storage.Unlock()

	userStatus.Lock()
	var ts = userStatus.m[update.Message.Chat.ID]
	slog.Debug(fmt.Sprintf("ChatId: %d", update.Message.Chat.ID))
	slog.Debug(fmt.Sprintf("ts before setPayment: %+v", ts))
	ts.setPayment(string(update.Message.Text))
	userStatus.m[update.Message.Chat.ID] = ts
	slog.Debug(fmt.Sprintf("ts after setPayment: %+v", ts))
	userStatus.Unlock()

	ConfirmationMsg(ctx, b, update)
}

// default button to remove the payment one

func defaultButton(ctx context.Context, b *bot.Bot, update *models.Update) {
	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Connection successful!",
	})
}

// Confirmation
func ConfirmationMsg(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	response, exists := userStatus.m[update.Message.Chat.ID]
	slog.Info(fmt.Sprintf("ChatId: %d confirmation response: %+v", update.Message.Chat.ID, response))
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Something went wrong with data!",
		})
	}
	var ts = userStatus.m[update.Message.Chat.ID]
	userName, _ := auth.GetUserName(update.Message.Chat.ID)
	db.Insert(userName, ts.sum, ts.category, ts.payment)

	var confirmationMsg = fmt.Sprintf("Amount:\t\t\t¥%d\nCategory:\t\t%s\nPayment:\t\t%s", response.sum, response.category, response.payment)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:			update.Message.Chat.ID,
		Text:			confirmationMsg,
		ReplyMarkup:	defaultReplyKeyboard,
	})

	userStatus.Lock()
	slog.Debug(fmt.Sprintf("ChatId: %d", update.Message.Chat.ID))
	slog.Debug(fmt.Sprintf("ts before reset: %+v", ts))
	ts.reset()
	userStatus.m[update.Message.Chat.ID] = ts
	slog.Debug(fmt.Sprintf("ts after reset: %+v", ts))
	userStatus.Unlock()

	mode.Storage.Lock()
	mode.Storage.M[update.Message.Chat.ID] = mode.InputAmount
	mode.Storage.Unlock()
}

// TODO: add confirmation step
func Confirmed(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Expence has confirmed",
	})
}
