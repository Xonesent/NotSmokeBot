package postgres

const (
	UsersTableName     = "user_schema.users"
	UsersInfoTableName = "user_schema.users_info"
	UserId             = "user_id"
	InfoId             = "info_id"
	TgId               = "tg_id"
	LastMes            = "last_mes"
	Nick               = "nick"
	Money              = "money"
	Progress           = "progress"
	CurrQuotation      = "curr_quotation"
)

var (
	InsertUserColumns = []string{
		TgId,
		LastMes,
		Progress,
		CurrQuotation,
	}
)

type CreateParams struct {
	TgId    int
	LastMes string
}

type UpdateParams struct {
	TgId    int
	LastMes string
}

type GetMesParams struct {
	TgId int
}

type UpdateNickParams struct {
	TgId int
	Nick string
}

type GetNickParams struct {
	TgId int
}

type UpdateMoneyParams struct {
	TgId  int
	Money int
}

type GetMoneyParams struct {
	TgId int
}

type UpdateQuotationParams struct {
	TgId      int
	Quotation int
}

type UpdateProgressParams struct {
	TgId     int
	Progress int
}
