package tg_buttons

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/error_handler"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel"
)

type ButtonHandler struct {
	cfg      *config.Config
	buttonUC ButtonUC
	errHDL   *error_handler.ErrorHandler
}

func NewButtonHandler(cfg *config.Config, buttonUC ButtonUC, errHDL *error_handler.ErrorHandler) *ButtonHandler {
	return &ButtonHandler{
		cfg:      cfg,
		buttonUC: buttonUC,
		errHDL:   errHDL,
	}
}

func (h *ButtonHandler) StartBot() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.StartBot")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.buttonUC.StartBot(ctx, sentMessageDTO)

		h.errHDL.HandleError(ctx, b, update, err)
	}
}
