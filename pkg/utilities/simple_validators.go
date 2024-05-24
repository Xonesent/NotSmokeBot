package utilities

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ValidateDefaultHandler() bot.MatchFunc {
	return func(update *models.Update) bool {
		if update.Message != nil && update.Message.Text != "/start" {
			return true
		}
		return false
	}
}
