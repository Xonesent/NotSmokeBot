package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	MongoId      primitive.ObjectID
	TgId         TgId
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
