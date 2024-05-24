package telegram

import (
	"NotSmokeBot/pkg/utilities"
	"github.com/go-telegram/bot"
)

func MapButtonRoutes(b *bot.Bot, h *ButtonHandler) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartBot())
	b.RegisterHandlerMatchFunc(utilities.ValidateDefaultHandler(), h.DefaultResponse())
}
