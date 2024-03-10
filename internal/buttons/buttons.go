package buttons

import (
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/postgres"
	"NotSmokeBot/pkg/utilities"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ShowProgress(ctx context.Context, b *bot.Bot, update *models.Update) {
	num, _ := postgres.GetProgress(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.Message.From.ID)})
	days := utilities.GetDays(num)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Вау, ты не курил %d %s подряд, Легенда!", num, days),
	})
}
