package error_handler

import (
	"NotSmokeBot/pkg/templates/errlst"
	"NotSmokeBot/pkg/templates/tg_resp"
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func ErrorHandler(ctx context.Context, b *bot.Bot, update *models.Update, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, errlst.AlreadyExists):
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: tg_resp.AlreadyExistResp})
	default:
		zap.L().Error(err.Error())
	}
}
