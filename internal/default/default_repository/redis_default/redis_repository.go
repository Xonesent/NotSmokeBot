package redis_default

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/model"
	"NotSmokeBot/pkg/dependences/tracer"
	"NotSmokeBot/pkg/errors/errlst"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"strconv"
	"time"
)

type DefaultRedisRepository struct {
	cfg *config.Config
	db  *redis.Client
}

func NewDefaultRedisRepository(cfg *config.Config, db *redis.Client) *DefaultRedisRepository {
	return &DefaultRedisRepository{
		cfg: cfg,
		db:  db,
	}
}

// PutUserInfoById Есть сомнения, что должно быть ключом, в итоге пришел к выводу с tg_id, так как у меня не будет кук,
// из которых я смогу доставать сложные сессионные ключи, а tg_id есть всегда и он уникален (ну и в названии ById)
func (r *DefaultRedisRepository) PutUserInfoById(ctx context.Context, userInfo *PutUserInfoById) error {
	ctx, span := otel.Tracer("").Start(ctx, "DefaultRedisRepository.PutUserInfoById")
	defer span.End()

	sessionBytes, err := json.Marshal(userInfo)
	if err != nil {
		return tracer.SpanSetErrWrap(span, errlst.ServerError, err, "DefaultRedisRepository.PutUserInfoById.Marshal")
	}

	key := strconv.Itoa(int(userInfo.TgId))

	_, err = r.db.Set(ctx, key, sessionBytes, time.Duration(r.cfg.Telegram.SessionTTL)*time.Second).Result()
	if err != nil {
		return tracer.SpanSetErrWrap(span, errlst.ServerError, err, "DefaultRedisRepository.PutUserInfoById.Set")
	}

	return nil
}

func (r *DefaultRedisRepository) GetUserInfoById(ctx context.Context, tgId model.TgId) (model.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DefaultRedisRepository.GetUserInfoById")
	defer span.End()

	user := model.User{}
	key := strconv.Itoa(int(tgId))

	valueString, err := r.db.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return user, tracer.SpanSetErrWrap(span, errlst.NothingFoundError, err, "DefaultRedisRepository.GetUserInfoById.Nil")
	} else if err != nil {
		return user, tracer.SpanSetErrWrap(span, errlst.ServerError, err, "DefaultRedisRepository.GetUserInfoById.Get")
	}

	err = json.Unmarshal([]byte(valueString), &user)
	if err != nil {
		return user, tracer.SpanSetErrWrap(span, errlst.ServerError, err, "DefaultRedisRepository.GetUserInfoById.Unmarshal")
	}

	return user, nil
}

func (r *DefaultRedisRepository) DelUserById(ctx context.Context, tgId model.TgId) error {
	ctx, span := otel.Tracer("").Start(ctx, "DefaultRedisRepository.DelUserById")
	defer span.End()

	key := strconv.Itoa(int(tgId))

	_, err := r.db.Del(ctx, key).Result()
	if err != nil {
		return tracer.SpanSetErrWrap(span, errlst.ServerError, err, "DefaultRedisRepository.DelUserById.Del")
	}

	return nil
}
