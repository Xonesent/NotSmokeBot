package buttons_repository

import (
	"NotSmokeBot/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SentMessage struct {
	Sender  model.TgId
	Message string
	ChatId  int64
}

type UpdateUserInfo struct {
	MongoId     primitive.ObjectID
	TgId        model.TgId
	LastMessage string
	Progress    *int64
	Quotation   *int64
	Nick        string
	Money       float32
	DeletedAt   bool
}

type User struct {
	MongoId     primitive.ObjectID `bson:"_id"`
	TgId        model.TgId         `bson:"tg_id"`
	LastMessage string             `bson:"last_mes"`
	Progress    int64              `bson:"progress"`
	Quotation   int64              `bson:"quotation"`
	Nick        string             `bson:"nick"`
	Money       float32            `bson:"money"`
	ChatId      int64              `bson:"chat_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	DeletedAt   time.Time          `bson:"deleted_at"`
}

func (d *User) toUserModel() model.User {
	return model.User{
		TgId:        d.TgId,
		LastMessage: d.LastMessage,
		Progress:    d.Progress,
		Quotation:   d.Quotation,
		Nick:        d.Nick,
		Money:       d.Money,
		ChatId:      d.ChatId,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
	}
}
