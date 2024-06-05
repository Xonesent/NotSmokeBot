package error_handler

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"NotSmokeBot/internal/model"
	"NotSmokeBot/pkg/errors/errlst"
	"NotSmokeBot/pkg/templates/tg_resp"
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	cfg            *config.Config
	defaultMNGRepo DefaultMNGRepo
	defaultRDRepo  DefaultRDRepo
}

func NewErrorHandler(cfg *config.Config, defaultMNGRepo DefaultMNGRepo, defaultRDRepo DefaultRDRepo) *ErrorHandler {
	return &ErrorHandler{
		cfg:            cfg,
		defaultMNGRepo: defaultMNGRepo,
		defaultRDRepo:  defaultRDRepo,
	}
}

// HandleError необходим не только для лога, но еще и используется для обработки ошибки (ответ пользователю, изменения полей бд и тд)
func (h *ErrorHandler) HandleError(ctx context.Context, b *bot.Bot, update *models.Update, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, errlst.AlreadyExistsError):
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: tg_resp.AlreadyExistResp})
	case errors.Is(err, bot.ErrorForbidden):
		err = h.defaultMNGRepo.UpdateUserByIds(ctx, mongo_default.UpdateUserInfo{TgId: model.TgId(update.Message.From.ID), DeletedAt: true})
		err = h.defaultRDRepo.DelUserById(ctx, model.TgId(update.Message.From.ID))
	}

	if err != nil {
		zap.L().Error(err.Error())
	}
}
