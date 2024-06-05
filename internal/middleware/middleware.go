package middleware

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/error_handler"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MDWManager struct {
	cfg            *config.Config
	defaultMNGRepo DefaultMNGRepo
	errHDL         *error_handler.ErrorHandler
}

func NewMiddleware(cfg *config.Config, defaultMNGRepo DefaultMNGRepo, errHDL *error_handler.ErrorHandler) *MDWManager {
	return &MDWManager{
		cfg:            cfg,
		defaultMNGRepo: defaultMNGRepo,
		errHDL:         errHDL,
	}
}

func (m *MDWManager) StartedMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx = context.WithValue(ctx, "bebra", "bebra")

		next(ctx, b, update)
	}
}
