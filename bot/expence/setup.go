package expence

import (
	"context"
	"fmt"

	"sync"
	"telegram-budget-bot/bot/auth"
	"telegram-budget-bot/bot/db"
	"telegram-budget-bot/bot/shared"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"
)

var defaultReplyKeyboard *reply.ReplyKeyboard
var categoryReplyKeyboard *reply.ReplyKeyboard
var paymentReplyKeyboard *reply.ReplyKeyboard

var userStatus = struct {
	sync.RWMutex
	m map[int64]transactionStatus
}{m: make(map[int64]transactionStatus)}

type transactionStatus struct {
	sum      int
	category string
	payment  string
}

func (ts *transactionStatus) setCategory(category string) {
	ts.category = category
}

func (ts *transactionStatus) setPayment(payment string) {
	ts.payment = payment
}

func (ts *transactionStatus) reset() {
	ts.sum = 0
	ts.category = ""
	ts.payment = ""
}

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
func InitCategoryKeyboard(b *bot.Bot) {
	categoryReplyKeyboard = reply.New(reply.WithPrefix("category_keyboard"), reply.IsSelective()).
		Button("🍽️ Food & 🧼 Household goods", b, bot.MatchTypeExact, ChoosenCategory).Button("🚇 Transport", b, bot.MatchTypeExact, ChoosenCategory).
		Row().
		Button("☕ Cafe & 🍻Parties", b, bot.MatchTypeExact, ChoosenCategory).Button("🐕 Dog", b, bot.MatchTypeExact, ChoosenCategory).
		Row().
		Button("🛒 Shopping", b, bot.MatchTypeExact, ChoosenCategory).Button("💧 Water & ⚡Electricity", b, bot.MatchTypeExact, ChoosenCategory).
		Row().
		Button("🏛️ Taxes", b, bot.MatchTypeExact, ChoosenCategory).Button("💄Beaty", b, bot.MatchTypeExact, ChoosenCategory).
		Row()
}

func ChooseCategoryMsg(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select example command from reply keyboard:",
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
	shared.Mode.Lock()
	shared.Mode.M[update.Message.Chat.ID] = "waiting_for_payment"
	shared.Mode.Unlock()

	userStatus.Lock()
	var ts = userStatus.m[update.Message.Chat.ID]
	
	fmt.Println("ts before categorySet: ", ts)
	ts.setCategory(string(update.Message.Text))
	userStatus.m[update.Message.Chat.ID] = ts
	fmt.Println("ts after categorySet: ", ts)
	userStatus.Unlock()

	ChoosePaymentMsg(ctx, b, update)
}

// Payment
func InitPaymentKeyboard(b *bot.Bot) {
	paymentReplyKeyboard = reply.New(reply.WithPrefix("payment_keyboard"), reply.IsOneTimeKeyboard()).
		Button("💵 Cash", b, bot.MatchTypeExact, ChoosenPayment).Button("💳 Credit Card", b, bot.MatchTypeExact, ChoosenPayment)
}

func ChoosePaymentMsg(ctx context.Context, b *bot.Bot, update *models.Update) {

	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select example command from reply keyboard:",
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

	shared.Mode.Lock()
	shared.Mode.M[update.Message.Chat.ID] = "waiting_for_confirmation"
	shared.Mode.Unlock()

	userStatus.Lock()
	var ts = userStatus.m[update.Message.Chat.ID]
	fmt.Println("ChatId: ", update.Message.Chat.ID)
	fmt.Println("ts before setPayment: ", ts)
	ts.setPayment(string(update.Message.Text))
	userStatus.m[update.Message.Chat.ID] = ts
	fmt.Println("ts after setPayment: ", ts)
	userStatus.Unlock()

	ConfirmationMsg(ctx, b, update)
}

// default button to remove the payment one
func InitDefaultKeyboard(b *bot.Bot) {
	defaultReplyKeyboard = reply.New(reply.WithPrefix("default_keyboard"), reply.IsOneTimeKeyboard()).
		Button("Dummy", b, bot.MatchTypeExact, defaultButton)
}

func defaultButton(ctx context.Context, b *bot.Bot, update *models.Update) {
	var _, ok = auth.GetUserName(update.Message.Chat.ID)
	if !ok {
		auth.SendMessageToUnregisteredUser(ctx, b, update)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "11!!!11!1!",
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
	fmt.Println("ChatId: ", update.Message.Chat.ID)
	fmt.Println("confirmation response: ", response)
	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Something went wrong with data!",
		})
	}
	var ts = userStatus.m[update.Message.Chat.ID]
	userName,_ := auth.GetUserName(update.Message.Chat.ID)
	db.Insert(userName, ts.sum, ts.category, ts.payment)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("%#v", response),
		ReplyMarkup: defaultReplyKeyboard,
	})

	userStatus.Lock()
	fmt.Println("ChatId: ", update.Message.Chat.ID)
	fmt.Println("ts before reset: ", ts)
	ts.reset()
	userStatus.m[update.Message.Chat.ID] = ts
	fmt.Println("ts after reset: ", ts)
	userStatus.Unlock()

	shared.Mode.Lock()
	shared.Mode.M[update.Message.Chat.ID] = ""
	shared.Mode.Unlock()
}

func Confirmed(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Expence has confirmed",
	})
}
