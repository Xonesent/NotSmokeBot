package main

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/beginning"
	"NotSmokeBot/internal/buttons"
	_default "NotSmokeBot/internal/default"
	"NotSmokeBot/internal/progress"
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/postgres"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	psqlDB, err := postgres.NewDB(cfg.Postgres)

	opts := []bot.Option{
		bot.WithDefaultHandler(_default.Echo),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, progress.Survey),
		bot.WithCallbackQueryDataHandler("quotation", bot.MatchTypePrefix, progress.QuotationNumber),
		bot.WithCallbackQueryDataHandler("begin", bot.MatchTypePrefix, beginning.Begin),
	}

	b, err := bot.New(os.Getenv("TG_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, beginning.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/progress", bot.MatchTypeExact, buttons.ShowProgress)

	server.Run(cfg, psqlDB, b)
}
