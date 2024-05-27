package error_handler

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/buttons/buttons_repository"
	"NotSmokeBot/internal/buttons/buttons_usecase"
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
	cfg           *config.Config
	buttonMNGRepo buttons_usecase.ButtonMNGRepo
}

func NewErrorHandler(cfg *config.Config, buttonMNGRepo buttons_usecase.ButtonMNGRepo) *ErrorHandler {
	return &ErrorHandler{
		cfg:           cfg,
		buttonMNGRepo: buttonMNGRepo,
	}
}

type ErrHandler interface {
	HandleError(ctx context.Context, b *bot.Bot, update *models.Update, err error)
}

func (h *ErrorHandler) HandleError(ctx context.Context, b *bot.Bot, update *models.Update, err error) {
	if err == nil {
		return
	}

	for err != nil {
		switch {
		case errors.Is(err, errlst.AlreadyExistsError):
			err = nil
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: tg_resp.AlreadyExistResp})
		case err.Error() == errlst.BlockErorr.Error():
			err = nil
			err = h.buttonMNGRepo.UpdateUserByIds(ctx, buttons_repository.UpdateUserInfo{TgId: model.TgId(update.Message.From.ID), DeletedAt: true})
		default:
			zap.L().Error(err.Error())
			err = nil
		}
	}
}
