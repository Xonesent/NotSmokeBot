package telegram

import (
	"NotSmokeBot/pkg/utils"
	"github.com/go-telegram/bot"
)

func MapButtonRoutes(b *bot.Bot, h *ButtonHandler) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartBot())
	b.RegisterHandlerMatchFunc(utils.ValidateDefaultHandler(), h.DefaultResponse())
}
