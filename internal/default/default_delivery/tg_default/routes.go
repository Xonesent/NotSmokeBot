package tg_default

import (
	"NotSmokeBot/internal/middleware"
	"NotSmokeBot/pkg/utils"
	"github.com/go-telegram/bot"
)

func MapDefaultRoutes(b *bot.Bot, h *DefaultHandler, mw *middleware.MDWManager) {
	b.RegisterHandlerMatchFunc(utils.ValidateDefaultHandler(), mw.StartedMiddleware(h.DefaultResponse()))
}
