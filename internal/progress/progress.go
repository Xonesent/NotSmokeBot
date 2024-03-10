package progress

import (
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/postgres"
	"NotSmokeBot/pkg/templates"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Survey(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	if update.CallbackQuery.Data == "button_1" {
		postgres.UpdateProgress(&server.Serve, ctx, postgres.UpdateProgressParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID), Progress: 0})
	} else if update.CallbackQuery.Data == "button_2" {
		num, _ := postgres.GetProgress(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
		postgres.UpdateProgress(&server.Serve, ctx, postgres.UpdateProgressParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID), Progress: num + 1})
	}

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}

func QuotationNumber(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	if update.CallbackQuery.Data == "quotation_1" {
		num, _ := postgres.GetQuotation(&server.Serve, ctx, postgres.GetMesParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID)})
		postgres.UpdateQuotation(&server.Serve, ctx, postgres.UpdateQuotationParams{TgId: int(update.CallbackQuery.Message.Message.Chat.ID), Quotation: (num + 1) % len(templates.QuotationsSlice)})
	}

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}
