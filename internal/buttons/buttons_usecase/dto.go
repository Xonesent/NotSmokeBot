package buttons_usecase

import (
	"NotSmokeBot/internal/buttons/buttons_repository"
	"NotSmokeBot/internal/model"
)

type StartMessage struct {
	Sender  model.TgId
	Message string
}

func (d *StartMessage) toStartMessage() buttons_repository.StartMessage {
	return buttons_repository.StartMessage{
		Sender:  d.Sender,
		Message: d.Message,
	}
}
