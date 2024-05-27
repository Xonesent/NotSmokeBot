package telegram

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/errors/error_handler"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel"
)

type ButtonHandler struct {
	cfg        *config.Config
	buttonUC   ButtonUC
	errHandler *error_handler.ErrorHandler
}

func NewButtonHandler(cfg *config.Config, buttonUC ButtonUC, errHandler *error_handler.ErrorHandler) *ButtonHandler {
	return &ButtonHandler{
		cfg:        cfg,
		buttonUC:   buttonUC,
		errHandler: errHandler,
	}
}

func (h *ButtonHandler) DefaultResponse() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.DefaultResponse")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.DefaultResponse(ctx, sentMessageDTO)

		h.errHandler.HandleError(ctx, b, update, err)
	}
}

func (h *ButtonHandler) StartBot() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.StartBot")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.StartBot(ctx, sentMessageDTO)

		h.errHandler.HandleError(ctx, b, update, err)
	}
}
