package default_usecase

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/dependences/tracer"
	"NotSmokeBot/pkg/errors/errlst"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-telegram/bot"
	"go.opentelemetry.io/otel"
)

type DefaultUseCase struct {
	cfg            *config.Config
	defaultMNGRepo DefaultMNGRepo
	trManager      *manager.Manager
	b              *bot.Bot
}

func NewDefaultUseCase(cfg *config.Config, defaultMNGRepo DefaultMNGRepo, trManager *manager.Manager, b *bot.Bot) *DefaultUseCase {
	return &DefaultUseCase{
		cfg:            cfg,
		defaultMNGRepo: defaultMNGRepo,
		trManager:      trManager,
		b:              b,
	}
}

func (u *DefaultUseCase) DefaultResponse(ctx context.Context, sentMessage SentMessage) error {
	ctx, span := otel.Tracer("").Start(ctx, "DefaultUseCase.DefaultResponse")
	defer span.End()

	if err := u.trManager.Do(ctx, func(ctx context.Context) error {
		users, err := u.defaultMNGRepo.FindUsersByFilter(ctx, sentMessage.toFindUserByFilter())
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return tracer.SpanSetErrWrap(span, errlst.NotStartedError, errlst.NotStartedError, "DefaultUseCase.DefaultResponse.len(users) == 0")
		} else if len(users) > 1 {
			return tracer.SpanSetErrWrap(span, errlst.NotStartedError, errlst.NotStartedError, "DefaultUseCase.DefaultResponse.len(users) == 0")
		}

		return nil
	}); err != nil {
		return err
	}

	updateLastMessageDTO := sentMessage.toUpdateLastMessage()

	err := u.defaultMNGRepo.UpdateUserByIds(ctx, updateLastMessageDTO)
	if err != nil {
		return err
	}

	return nil
}
