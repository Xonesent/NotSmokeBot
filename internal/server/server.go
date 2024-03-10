package server

import (
	"NotSmokeBot/config"
	"context"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
)

type Database struct {
	Db        *sqlx.DB
	TxGetter  *trmsqlx.CtxGetter
	TrManager *manager.Manager
}

type Server struct {
	Cfg  *config.Config
	PgDB Database
	Bot  *bot.Bot
}

func NewServer(
	cfg *config.Config,
	pgDB *sqlx.DB,
	bot *bot.Bot,
) Server {
	return Server{
		Cfg: cfg,
		PgDB: Database{
			Db:        pgDB,
			TxGetter:  trmsqlx.DefaultCtxGetter,
			TrManager: manager.Must(trmsqlx.NewDefaultFactory(pgDB)),
		},
		Bot: bot,
	}
}

var Serve Server

func Run(cfg *config.Config, pgDB *sqlx.DB, bot *bot.Bot) {
	Serve = NewServer(cfg, pgDB, bot)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	Serve.Bot.Start(ctx)
}
