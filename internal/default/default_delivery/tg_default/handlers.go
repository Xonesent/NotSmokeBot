package tg_default

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/error_handler"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel"
)

type DefaultHandler struct {
	cfg       *config.Config
	defaultUC DefaultUC
	errHDL    *error_handler.ErrorHandler
}

func NewDefaultHandler(cfg *config.Config, defaultUC DefaultUC, errHDL *error_handler.ErrorHandler) *DefaultHandler {
	return &DefaultHandler{
		cfg:       cfg,
		defaultUC: defaultUC,
		errHDL:    errHDL,
	}
}

func (h *DefaultHandler) DefaultResponse() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "DefaultHandler.DefaultResponse")
		defer span.End()

		sentMessageDTO := toSentMessage(update)

		err := h.defaultUC.DefaultResponse(ctx, sentMessageDTO)

		h.errHDL.HandleError(ctx, b, update, err)
	}
}
