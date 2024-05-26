package telegram

import (
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"context"
	"github.com/go-telegram/bot"
)

type ButtonUC interface {
	DefaultResponse(ctx context.Context, sentMessage buttons_usecase.SentMessage) error
	StartBot(ctx context.Context, sentMessage buttons_usecase.SentMessage) error
}

type ButtonHDL interface {
	DefaultResponse() bot.HandlerFunc
	StartBot() bot.HandlerFunc
}
