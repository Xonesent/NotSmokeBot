package buttons_usecase

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/templates/tg_resp"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-telegram/bot"
	"go.opentelemetry.io/otel"
)

type ButtonUseCase struct {
	cfg           *config.Config
	buttonMNGRepo ButtonMNGRepo
	trManager     *manager.Manager
	b             *bot.Bot
}

func NewButtonUseCase(cfg *config.Config, buttonMNGRepo ButtonMNGRepo, trManager *manager.Manager, b *bot.Bot) *ButtonUseCase {
	return &ButtonUseCase{
		cfg:           cfg,
		buttonMNGRepo: buttonMNGRepo,
		trManager:     trManager,
		b:             b,
	}
}

func (u *ButtonUseCase) StartBot(ctx context.Context, sentMessage SentMessage) error {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonUseCase.StartBot")
	defer span.End()

	if err := u.trManager.Do(ctx, func(ctx context.Context) error {
		_, err := u.buttonMNGRepo.InsertNewUser(ctx, sentMessage.toStartMessage())
		if err != nil {
			return err
		}

		//time.Sleep(7 * time.Second)

		if _, err = u.b.SendMessage(ctx, &bot.SendMessageParams{ChatID: sentMessage.ChatId, Text: tg_resp.StartResp}); err != nil {
			return err
		}
		if _, err = u.b.SendMessage(ctx, &bot.SendMessageParams{ChatID: sentMessage.ChatId, Text: tg_resp.RegistrationReq}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
