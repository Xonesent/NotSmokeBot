package telegram

import (
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"context"
	"github.com/go-telegram/bot"
)

type ButtonUC interface {
	StartBot(ctx context.Context, startMessage buttons_usecase.StartMessage) error
}

type ButtonHDL interface {
	DefaultResponse() bot.HandlerFunc
	StartBot() bot.HandlerFunc
}
