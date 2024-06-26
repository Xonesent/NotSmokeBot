package server

import (
	"NotSmokeBot/config"
	"context"
	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg       *config.Config
	mngClient *mongo.Client
	redis     *redis.Client
	bot       *bot.Bot
}

func NewServer(
	cfg *config.Config,
	mngClient *mongo.Client,
	redis *redis.Client,
	bot *bot.Bot,
) *Server {
	return &Server{
		cfg:       cfg,
		mngClient: mngClient,
		redis:     redis,
		bot:       bot,
	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	zap.L().Info("Trying to run server...")

	err := s.MapHandlers()
	if err != nil {
		zap.L().Fatal("Error MapHandlers running server", zap.Error(err))
		return err
	}

	go func() {
		s.bot.Start(ctx)
	}()

	zap.L().Info("Server is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	zap.L().Info("Server is closing")

	return nil
}
