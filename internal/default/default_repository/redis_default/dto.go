package redis_default

import (
	"NotSmokeBot/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PutUserInfoById struct {
	MongoId      primitive.ObjectID
	TgId         model.TgId
	LastMessage  string
	Progress     int64
	Quotation    int64
	ActionStatus string
	Nick         string
	Money        float32
	ChatId       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
