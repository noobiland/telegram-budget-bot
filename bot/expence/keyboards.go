package expence

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/reply"
)

var defaultReplyKeyboard *reply.ReplyKeyboard
var categoryReplyKeyboard *reply.ReplyKeyboard
var paymentReplyKeyboard *reply.ReplyKeyboard

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

func InitPaymentKeyboard(b *bot.Bot) {
	paymentReplyKeyboard = reply.New(reply.WithPrefix("payment_keyboard"), reply.IsOneTimeKeyboard()).
		Button("💵 Cash", b, bot.MatchTypeExact, ChoosenPayment).Button("💳 Credit Card", b, bot.MatchTypeExact, ChoosenPayment)
}

func InitDefaultKeyboard(b *bot.Bot) {
	defaultReplyKeyboard = reply.New(reply.WithPrefix("default_keyboard")).
		Button("Check Connection...", b, bot.MatchTypeExact, defaultButton)
}