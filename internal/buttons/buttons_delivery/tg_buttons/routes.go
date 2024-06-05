package tg_buttons

import (
	"NotSmokeBot/internal/middleware"
	"github.com/go-telegram/bot"
)

func MapButtonRoutes(b *bot.Bot, h *ButtonHandler, mw *middleware.MDWManager) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, mw.StartedMiddleware(h.StartBot()))
}
