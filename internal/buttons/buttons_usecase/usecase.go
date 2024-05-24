package buttons_usecase

import (
	"NotSmokeBot/config"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.opentelemetry.io/otel"
)

type ButtonUseCase struct {
	buttonMNGRepo ButtonMNGRepo
	cfg           *config.Config
	trManager     *manager.Manager
}

func NewButtonUseCase(buttonMNGRepo ButtonMNGRepo, cfg *config.Config, trManager *manager.Manager) *ButtonUseCase {
	return &ButtonUseCase{
		buttonMNGRepo: buttonMNGRepo,
		cfg:           cfg,
		trManager:     trManager,
	}
}

func (u *ButtonUseCase) StartBot(ctx context.Context, startMessage StartMessage) error {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonUseCase.StartBot")
	defer span.End()

	if err := u.trManager.Do(ctx, func(ctx context.Context) error {
		_, err := u.buttonMNGRepo.InsertNewUser(ctx, startMessage.toStartMessage())
		_, err = u.buttonMNGRepo.InsertNewUser(ctx, startMessage.toStartMessage())
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
