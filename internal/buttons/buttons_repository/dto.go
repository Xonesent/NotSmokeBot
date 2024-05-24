package buttons_repository

import (
	"NotSmokeBot/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StartMessage struct {
	Sender  model.TgId
	Message string
}

type User struct {
	MongoId     primitive.ObjectID `bson:"_id"`
	TgId        model.TgId         `bson:"tg_id"`
	LastMessage string             `bson:"last_mes"`
	Progress    int64              `bson:"progress"`
	Quotation   int64              `bson:"quotation"`
	Nick        string             `bson:"nick"`
	Money       float32            `bson:"money"`
}

func (d *User) toUserModel() model.User {
	return model.User{
		TgId:        d.TgId,
		LastMessage: d.LastMessage,
		Progress:    d.Progress,
		Quotation:   d.Quotation,
		Nick:        d.Nick,
		Money:       d.Money,
	}
}
