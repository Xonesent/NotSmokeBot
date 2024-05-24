package main

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/dependences/mongo"
	"NotSmokeBot/pkg/tools/logger"
	"context"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log"
)

func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatalf("Error to init logger: %v\n", err)
	}

	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("Error loading env variables", zap.Error(err))
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Error loading config", zap.Error(err))
	}

	//tracer.NewTracer(cfg)
	//defer func(tracer *tracesdk.TracerProvider) {
	//	if err := tracer.Shutdown(ctx); err != nil {
	//		zap.L().Error("Jaeger shutdown error", zap.Error(err))
	//	} else {
	//		zap.L().Info("Jaeger shutdown properly")
	//	}
	//}(constant.Tracer)

	mngClient, err := mongodb.NewDB(cfg)
	if err != nil {
		zap.L().Fatal("Error connecting mongoDB", zap.Error(err))
	}
	defer func(mngClient *mongo.Client) {
		if err = mngClient.Disconnect(context.Background()); err != nil {
			zap.L().Error("MongoDB disconnect error", zap.Error(err))
		} else {
			zap.L().Info("MongoDB closed properly")
		}
	}(mngClient)

	opts := []bot.Option{}
	b, err := bot.New(cfg.Telegram.Token, opts...)

	s := server.NewServer(
		cfg,
		mngClient,
		b,
	)
	if err = s.Run(); err != nil {
		zap.L().Fatal("Cannot start server", zap.Error(err))
	}
}
