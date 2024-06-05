package tg_buttons

import (
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"NotSmokeBot/internal/model"
	"github.com/go-telegram/bot/models"
)

func toSentMessage(update *models.Update) buttons_usecase.SentMessage {
	return buttons_usecase.SentMessage{
		Sender:  model.TgId(update.Message.From.ID),
		ChatId:  update.Message.Chat.ID,
		Message: update.Message.Text,
	}
}
