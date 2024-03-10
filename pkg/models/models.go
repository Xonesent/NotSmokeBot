package dto

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Db       *sqlx.DB
	TxGetter *trmsqlx.CtxGetter
}
