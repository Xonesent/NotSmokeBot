package main

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/server"
	"NotSmokeBot/pkg/dependences/mongo"
	rediss "NotSmokeBot/pkg/dependences/redis"
	"NotSmokeBot/pkg/dependences/tracer"
	"NotSmokeBot/pkg/tools/logger"
	"context"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := context.Background()

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

	tracerJaeger := tracer.NewTracer(cfg)
	defer func(tracerJaeger *tracesdk.TracerProvider) {
		if err := tracerJaeger.Shutdown(ctx); err != nil {
			zap.L().Error("Jaeger shutdown error", zap.Error(err))
		} else {
			zap.L().Info("Jaeger shutdown properly")
		}
	}(tracerJaeger)

	mngClient, err := mongodb.NewMongoDB(cfg)
	if err != nil {
		zap.L().Fatal("Error connecting mongoDB", zap.Error(err))
	}
	defer func(mngClient *mongo.Client) {
		if err = mngClient.Disconnect(ctx); err != nil {
			zap.L().Error("MongoDB disconnect error", zap.Error(err))
		} else {
			zap.L().Info("MongoDB closed properly")
		}
	}(mngClient)

	redisClient, err := rediss.NewRedisClient(cfg)
	if err != nil {
		zap.L().Fatal("Error connecting redis", zap.Error(err))
	}
	defer func(redisClient *redis.Client) {
		if err = redisClient.Close(); err != nil {
			zap.L().Error("Redis close error", zap.Error(err))
		} else {
			zap.L().Info("Redis closed properly")
		}
	}(redisClient)

	opts := []bot.Option{}
	b, err := bot.New(cfg.Telegram.Token, opts...)

	s := server.NewServer(
		cfg,
		mngClient,
		redisClient,
		b,
	)
	if err = s.Run(); err != nil {
		zap.L().Fatal("Cannot start server", zap.Error(err))
	}
}
