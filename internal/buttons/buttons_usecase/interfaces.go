package buttons_usecase

import (
	"NotSmokeBot/internal/buttons/buttons_repository"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ButtonMNGRepo interface {
	InsertNewUser(ctx context.Context, startMessage buttons_repository.StartMessage) (primitive.ObjectID, error)
}
