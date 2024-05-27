package server

import (
	"NotSmokeBot/internal/buttons/buttons_delivery/telegram"
	"NotSmokeBot/internal/buttons/buttons_repository"
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"NotSmokeBot/pkg/errors/error_handler"
	trmmongo "github.com/avito-tech/go-transaction-manager/drivers/mongo/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func (s *Server) MapHandlers() error {
	buttonMNGRepository := buttons_repository.NewButtonsMNGRepository(s.cfg, s.mngClient.Database(s.cfg.Mongo.Database))

	trManager := manager.Must(
		trmmongo.NewDefaultFactory(s.mngClient),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)

	buttonUseCase := buttons_usecase.NewButtonUseCase(s.cfg, buttonMNGRepository, trManager, s.bot)

	errHandler := error_handler.NewErrorHandler(s.cfg, buttonMNGRepository)

	buttonHandlers := telegram.NewButtonHandler(s.cfg, buttonUseCase, errHandler)

	telegram.MapButtonRoutes(s.bot, buttonHandlers)
	return nil
}
