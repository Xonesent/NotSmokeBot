package middleware

import (
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"context"
)

type DefaultMNGRepo interface {
	UpdateUserByIds(ctx context.Context, updateUserInfo mongo_default.UpdateUserInfo) error
}
