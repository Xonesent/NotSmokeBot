package error_handler

import (
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"NotSmokeBot/internal/model"
	"context"
)

type DefaultMNGRepo interface {
	UpdateUserByIds(ctx context.Context, updateUserInfo mongo_default.UpdateUserInfo) error
}

type DefaultRDRepo interface {
	DelUserById(ctx context.Context, tgId model.TgId) error
}
