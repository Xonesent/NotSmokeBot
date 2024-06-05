package server

import (
	"NotSmokeBot/internal/buttons/buttons_delivery/tg_buttons"
	"NotSmokeBot/internal/buttons/buttons_repository"
	"NotSmokeBot/internal/buttons/buttons_usecase"
	"NotSmokeBot/internal/default/default_delivery/tg_default"
	"NotSmokeBot/internal/default/default_repository/mongo_default"
	"NotSmokeBot/internal/default/default_repository/redis_default"
	"NotSmokeBot/internal/default/default_usecase"
	"NotSmokeBot/internal/error_handler"
	"NotSmokeBot/internal/middleware"
	trmmongo "github.com/avito-tech/go-transaction-manager/drivers/mongo/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func (s *Server) MapHandlers() error {
	buttonMNGRepository := buttons_repository.NewButtonsMNGRepository(s.cfg, s.mngClient.Database(s.cfg.Mongo.Database))
	defaultMNGRepository := mongo_default.NewDefaultMNGRepository(s.cfg, s.mngClient.Database(s.cfg.Mongo.Database))
	defaultRedisRepository := redis_default.NewDefaultRedisRepository(s.cfg, s.redis)

	trManager := manager.Must(
		trmmongo.NewDefaultFactory(s.mngClient),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)

	buttonUseCase := buttons_usecase.NewButtonUseCase(s.cfg, buttonMNGRepository, trManager, s.bot)
	defaultUseCase := default_usecase.NewDefaultUseCase(s.cfg, defaultMNGRepository, trManager, s.bot)

	errHandler := error_handler.NewErrorHandler(s.cfg, defaultMNGRepository, defaultRedisRepository)

	buttonHandlers := tg_buttons.NewButtonHandler(s.cfg, buttonUseCase, errHandler)
	defaultHandlers := tg_default.NewDefaultHandler(s.cfg, defaultUseCase, errHandler)

	middlewareManager := middleware.NewMiddleware(s.cfg, defaultMNGRepository, errHandler)

	tg_buttons.MapButtonRoutes(s.bot, buttonHandlers, middlewareManager)
	tg_default.MapDefaultRoutes(s.bot, defaultHandlers, middlewareManager)
	return nil
}
