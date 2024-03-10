package postgres

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/server"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(cfg config.ConfigPg) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(r *server.Server, ctx context.Context, UserParams CreateParams) (int, error) {
	query, args, err := sq.Insert(UsersTableName).
		Columns(InsertUserColumns...).
		Values(
			UserParams.TgId,
			UserParams.LastMes,
			0,
			0,
		).
		Suffix("RETURNING " + UserId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var userID int
	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	err = tr.QueryRow(query, args...).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func CreateInfo(r *server.Server, ctx context.Context, UserParams int) (int, error) {
	query, args, err := sq.Insert(UsersInfoTableName).
		Columns(UserId).
		Values(UserParams).
		Suffix("RETURNING " + InfoId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var infoID int
	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	err = tr.QueryRow(query, args...).Scan(&infoID)
	if err != nil {
		return -1, err
	}

	return infoID, nil
}

func UpdateMessage(r *server.Server, ctx context.Context, UserParams UpdateParams) error {
	query, args, err := sq.Update(UsersTableName).
		Set(LastMes, UserParams.LastMes).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetMessage(r *server.Server, ctx context.Context, UserParams GetMesParams) (string, error) {
	query, args, err := sq.Select(LastMes).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var lastMes string
	err = r.PgDB.Db.QueryRowContext(ctx, query, args...).Scan(&lastMes)
	if err != nil {
		return "", err
	}

	return lastMes, nil
}

func UpdateNick(r *server.Server, ctx context.Context, UserParams UpdateNickParams) error {
	//query, args, err := sq.Update(UsersInfoTableName).
	//	Set(Nick, UserParams.Nick).
	//	Where(sq.Eq{
	//		UserId: sq.Expr(sq.Select(UserId).
	//			From(UsersTableName).
	//			Where(sq.Eq{TgId: UserParams.TgId}).
	//			PlaceholderFormat(sq.Dollar).
	//			ToSql(),
	//		)}).
	//	PlaceholderFormat(sq.Dollar).
	//	ToSql()
	//if err != nil {
	//	return err
	//}

	subQuery, args, err := sq.Select(UserId).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var userId int
	err = r.PgDB.Db.QueryRowContext(ctx, subQuery, args...).Scan(&userId)
	if err != nil {
		return err
	}

	query, args, err := sq.Update(UsersInfoTableName).
		Set(Nick, UserParams.Nick).
		Where(sq.Eq{UserId: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetNick(r *server.Server, ctx context.Context, UserParams GetNickParams) (string, error) {
	subQuery, args, err := sq.Select(UserId).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var userId int
	err = r.PgDB.Db.QueryRowContext(ctx, subQuery, args...).Scan(&userId)
	if err != nil {
		return "", err
	}

	query, args, err := sq.Select(Nick).
		From(UsersInfoTableName).
		Where(sq.Eq{UserId: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var Nickname string
	err = r.PgDB.Db.QueryRowContext(ctx, query, args...).Scan(&Nickname)
	if err != nil {
		return "", err
	}

	return Nickname, err
}

func UpdateMoney(r *server.Server, ctx context.Context, UserParams UpdateMoneyParams) error {
	subQuery, args, err := sq.Select(UserId).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var userId int
	err = r.PgDB.Db.QueryRowContext(ctx, subQuery, args...).Scan(&userId)
	if err != nil {
		return err
	}

	query, args, err := sq.Update(UsersInfoTableName).
		Set(Money, UserParams.Money).
		Where(sq.Eq{UserId: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetMoney(r *server.Server, ctx context.Context, UserParams GetMoneyParams) (int, error) {
	subQuery, args, err := sq.Select(UserId).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var userId int
	err = r.PgDB.Db.QueryRowContext(ctx, subQuery, args...).Scan(&userId)
	if err != nil {
		return -1, err
	}

	query, args, err := sq.Select(Money).
		From(UsersInfoTableName).
		Where(sq.Eq{UserId: userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var Money1 int
	err = r.PgDB.Db.QueryRowContext(ctx, query, args...).Scan(&Money1)
	if err != nil {
		return -1, err
	}

	return Money1, err
}

func GetQuotation(r *server.Server, ctx context.Context, UserParams GetMesParams) (int, error) {
	query, args, err := sq.Select(CurrQuotation).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var num int
	err = r.PgDB.Db.QueryRowContext(ctx, query, args...).Scan(&num)
	if err != nil {
		return -1, err
	}

	return num, nil
}

func UpdateQuotation(r *server.Server, ctx context.Context, UserParams UpdateQuotationParams) error {
	query, args, err := sq.Update(UsersTableName).
		Set(CurrQuotation, UserParams.Quotation).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetProgress(r *server.Server, ctx context.Context, UserParams GetMesParams) (int, error) {
	query, args, err := sq.Select(Progress).
		From(UsersTableName).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, err
	}

	var num int
	err = r.PgDB.Db.QueryRowContext(ctx, query, args...).Scan(&num)
	if err != nil {
		return -1, err
	}

	return num, nil
}

func UpdateProgress(r *server.Server, ctx context.Context, UserParams UpdateProgressParams) error {
	query, args, err := sq.Update(UsersTableName).
		Set(Progress, UserParams.Progress).
		Where(sq.Eq{TgId: UserParams.TgId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	tr := r.PgDB.TxGetter.DefaultTrOrDB(ctx, r.PgDB.Db)
	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
