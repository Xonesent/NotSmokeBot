package buttons_usecase

import (
	"NotSmokeBot/config"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-telegram/bot"
	"go.opentelemetry.io/otel"
)

type ButtonUseCase struct {
	buttonMNGRepo ButtonMNGRepo
	cfg           *config.Config
	trManager     *manager.Manager
	b             *bot.Bot
}

func NewButtonUseCase(buttonMNGRepo ButtonMNGRepo, cfg *config.Config, trManager *manager.Manager, b *bot.Bot) *ButtonUseCase {
	return &ButtonUseCase{
		buttonMNGRepo: buttonMNGRepo,
		cfg:           cfg,
		trManager:     trManager,
		b:             b,
	}
}

func (u *ButtonUseCase) DefaultResponse(ctx context.Context, sentMessage SentMessage) error {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonUseCase.DefaultResponse")
	defer span.End()

	updateLastMessageDTO := sentMessage.toUpdateLastMessage()

	err := u.buttonMNGRepo.UpdateUserByIds(ctx, updateLastMessageDTO)
	if err != nil {
		return err
	}

	return nil
}

func (u *ButtonUseCase) StartBot(ctx context.Context, sentMessage SentMessage) error {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonUseCase.StartBot")
	defer span.End()

	if err := u.trManager.Do(ctx, func(ctx context.Context) error {
		_, err := u.buttonMNGRepo.InsertNewUser(ctx, sentMessage.toStartMessage())
		if err != nil {
			return err
		}
		//if _, err := u.b.SendMessage(ctx, &bot.SendMessageParams{ChatID: sentMessage.ChatId, Text: tg_resp.AlreadyExistResp})

		return nil
	}); err != nil {
		return err
	}

	return nil
}
