package tg_buttons

import (
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"context"
	"github.com/go-telegram/bot"
)

type ButtonUC interface {
	StartBot(ctx context.Context, sentMessage buttons_usecase.SentMessage) error
}

type ButtonHDL interface {
	StartBot() bot.HandlerFunc
}
