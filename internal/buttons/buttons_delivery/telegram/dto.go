package telegram

import (
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"NotSmokeBot/internal/model"
	"github.com/go-telegram/bot/models"
)

func toStartMessage(update *models.Update) buttons_usecase.StartMessage {
	return buttons_usecase.StartMessage{
		Sender:  model.TgId(update.Message.From.ID),
		Message: update.Message.Text,
	}
}
