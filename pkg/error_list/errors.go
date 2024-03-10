package err_list

const (
	UserExist        = "pq: duplicate key value violates unique constraint \"users_tg_id_key\""
	TransactionExist = "pq: current transaction is aborted, commands ignored until end of transaction block"
)
