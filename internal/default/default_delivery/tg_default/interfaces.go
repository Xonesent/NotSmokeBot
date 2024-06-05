package tg_default

import (
	"NotSmokeBot/internal/default/default_usecase"
	"NotSmokeBot/internal/error_handler"
	"context"
	"github.com/go-telegram/bot"
)

type DefaultUC interface {
	DefaultResponse(ctx context.Context, sentMessage default_usecase.SentMessage) error
}

type DefaultHDL interface {
	DefaultResponse(errHDL *error_handler.ErrorHandler) bot.HandlerFunc
}
