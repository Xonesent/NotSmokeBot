package buttons_usecase

import (
	"NotSmokeBot/internal/buttons/buttons_repository"
	"NotSmokeBot/internal/model"
)

type SentMessage struct {
	Sender  model.TgId
	Message string
	ChatId  int64
}

func (d *SentMessage) toStartMessage() buttons_repository.SentMessage {
	return buttons_repository.SentMessage{
		Sender:  d.Sender,
		Message: d.Message,
		ChatId:  d.ChatId,
	}
}

func (d *SentMessage) toUpdateLastMessage() buttons_repository.UpdateUserInfo {
	return buttons_repository.UpdateUserInfo{
		TgId:        d.Sender,
		LastMessage: d.Message,
	}
}
