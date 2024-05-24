package telegram

import (
	"NotSmokeBot/config"
	"context"
	"fmt"
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
		mes, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "DEFAULT",
		})

		if err != nil {
			zap.L().Info(err.Error())
		} else {
			zap.L().Info(fmt.Sprintf("%v,", mes))
		}
	}
}

func (h *ButtonHandler) StartBot() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := otel.Tracer("").Start(ctx, "ButtonHandler.StartBot")
		defer span.End()

		startMessageDTO := toStartMessage(update)

		h.buttonUC.StartBot(ctx, startMessageDTO)
	}
}
