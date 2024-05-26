package telegram

import (
	"NotSmokeBot/config"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type ButtonHandler struct {
	buttonUC ButtonUC
	cfg      *config.Config
}

func NewButtonHandler(buttonUC ButtonUC, cfg *config.Config) *ButtonHandler {
	return &ButtonHandler{
		buttonUC: buttonUC,
		cfg:      cfg,
	}
}

func (h *ButtonHandler) DefaultResponse() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.DefaultResponse")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.DefaultResponse(ctx, sentMessageDTO)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}

func (h *ButtonHandler) StartBot() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.StartBot")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.StartBot(ctx, sentMessageDTO)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}
