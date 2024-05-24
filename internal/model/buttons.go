package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	MongoId     primitive.ObjectID
	TgId        TgId
	LastMessage string
	Progress    int64
	Quotation   int64
	Nick        string
	Money       float32
}
