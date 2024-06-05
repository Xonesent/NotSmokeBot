package default_usecase

import (
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"NotSmokeBot/internal/model"
)

type SentMessage struct {
	Sender  model.TgId
	Message string
	ChatId  int64
}

func (d *SentMessage) toStartMessage() mongo_default.SentMessage {
	return mongo_default.SentMessage{
		Sender:  d.Sender,
		Message: d.Message,
		ChatId:  d.ChatId,
	}
}

func (d *SentMessage) toUpdateLastMessage() mongo_default.UpdateUserInfo {
	return mongo_default.UpdateUserInfo{
		TgId:        d.Sender,
		LastMessage: d.Message,
	}
}

func (d *SentMessage) toFindUserByFilter() mongo_default.FindUsersByFilter {
	return mongo_default.FindUsersByFilter{
		TgId: []model.TgId{d.Sender},
	}
}
