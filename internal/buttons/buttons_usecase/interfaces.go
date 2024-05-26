package buttons_usecase

import (
	"NotSmokeBot/internal/buttons/buttons_repository"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ButtonMNGRepo interface {
	InsertNewUser(ctx context.Context, sentMessage buttons_repository.SentMessage) (primitive.ObjectID, error)
	UpdateUserByIds(ctx context.Context, updateUserInfo buttons_repository.UpdateUserInfo) error
}
