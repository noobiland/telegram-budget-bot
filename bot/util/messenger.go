package util

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Messenger interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}
