package telegram

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/error_handler"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel"
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

		error_handler.ErrorHandler(ctx, b, update, err)
	}
}

func (h *ButtonHandler) StartBot() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.StartBot")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.StartBot(ctx, sentMessageDTO)

		error_handler.ErrorHandler(ctx, b, update, err)
	}
}
