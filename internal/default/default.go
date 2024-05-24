package _default

import (
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/postgres"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

func Echo(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := postgres.UpdateMessage(&server.Serve, ctx, postgres.UpdateParams{TgId: int(update.Message.From.ID), LastMes: update.Message.Text})
	if err != nil {
		log.Print(err.Error())
	}
	//b.DeleteMessage(ctx, &bot.DeleteMessageParams{
	//	ChatID:    update.Message.Chat.ID,
	//	MessageID: update.Message.ID,
	//})
}
