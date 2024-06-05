package default_usecase

import (
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"NotSmokeBot/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultMNGRepo interface {
	InsertNewUser(ctx context.Context, sentMessage mongo_default.SentMessage) (primitive.ObjectID, error)
	UpdateUserByIds(ctx context.Context, updateUserInfo mongo_default.UpdateUserInfo) error
	FindUsersByFilter(ctx context.Context, findUsersByFilter mongo_default.FindUsersByFilter) ([]model.User, error)
}

type DefaultRDRepo interface {
}
