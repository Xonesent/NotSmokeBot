package tg_default

import (
	"NotSmokeBot/internal/default/default_usecase"
	"NotSmokeBot/internal/model"
	"github.com/go-telegram/bot/models"
)

func toSentMessage(update *models.Update) default_usecase.SentMessage {
	return default_usecase.SentMessage{
		Sender:  model.TgId(update.Message.From.ID),
		ChatId:  update.Message.Chat.ID,
		Message: update.Message.Text,
	}
}
