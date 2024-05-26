package buttons_repository

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/store/bson_queries"
	"NotSmokeBot/pkg/dependences/tracer"
	"NotSmokeBot/pkg/templates/errlst"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"time"
)

type ButtonsMNGRepository struct {
	cfg *config.Config
	db  *mongo.Database
}

func NewButtonsMNGRepository(cfg *config.Config, db *mongo.Database) *ButtonsMNGRepository {
	return &ButtonsMNGRepository{
		cfg: cfg,
		db:  db,
	}
}

func (r *ButtonsMNGRepository) InsertNewUser(ctx context.Context, sentMessage SentMessage) (primitive.ObjectID, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonsMNGRepository.InsertNewUser")
	defer span.End()

	collection := r.db.Collection(bson_queries.UsersCollection)
	insertResult, err := collection.InsertOne(context.TODO(),
		bson.M{
			bson_queries.TgIdColumnName:        sentMessage.Sender,
			bson_queries.LastMessageColumnName: sentMessage.Message,
			bson_queries.ProgressColumnName:    0,
			bson_queries.QuotationColumnName:   0,
			bson_queries.NickColumnName:        "",
			bson_queries.MoneyColumnName:       0.0,
			bson_queries.ChatIdColumnName:      sentMessage.ChatId,
			bson_queries.CreatedAtColumnName:   time.Now(),
			bson_queries.UpdatedAtColumnName:   time.Now(),
			bson_queries.DeletedAtColumnName:   nil,
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.AlreadyExists, err, "ButtonsMNGRepository.InsertNewUser.IsDuplicateKeyError")
		}
		return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.ServerError, err, "ButtonsMNGRepository.InsertNewUser.InsertOne")
	}

	value, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.ConvertionError, errors.New("primitive.ObjectID convert error"), "ButtonsMNGRepository.InsertNewUser.ok")
	}

	return value, nil
}

func (r *ButtonsMNGRepository) UpdateUserByIds(ctx context.Context, updateUserInfo UpdateUserInfo) error {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonsMNGRepository.UpdateUserByIds")
	defer span.End()

	collection := r.db.Collection(bson_queries.UsersCollection)
	filter, update, err := getUpdateParams(updateUserInfo)
	if err != nil {
		return tracer.SpanSetErrWrap(span, errlst.NilUpdateFields, err, "ButtonsMNGRepository.UpdateUserByIds.getUpdateParams")
	}

	_, err = collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return tracer.SpanSetErrWrap(span, errlst.ServerError, err, "ButtonsMNGRepository.UpdateUserByIds.UpdateMany")
	}

	return nil
}

func getUpdateParams(updateUserInfo UpdateUserInfo) (bson.M, bson.M, error) {
	filter := bson.M{}
	if updateUserInfo.MongoId != [12]byte{} {
		filter[bson_queries.IdColumnName] = updateUserInfo.MongoId
	}
	if updateUserInfo.TgId != 0 {
		filter[bson_queries.TgIdColumnName] = updateUserInfo.TgId
	}

	update := bson.M{}
	if updateUserInfo.LastMessage != "" {
		update[bson_queries.LastMessageColumnName] = updateUserInfo.LastMessage
	}
	if updateUserInfo.Progress != nil {
		update[bson_queries.ProgressColumnName] = *updateUserInfo.Progress
	}
	if updateUserInfo.Quotation != nil {
		update[bson_queries.QuotationColumnName] = *updateUserInfo.Quotation
	}
	if updateUserInfo.Nick != "" {
		update[bson_queries.NickColumnName] = updateUserInfo.Nick
	}
	if updateUserInfo.Money != 0.0 {
		update[bson_queries.MoneyColumnName] = updateUserInfo.Money
	}
	if updateUserInfo.DeletedAt != false {
		update[bson_queries.DeletedAtColumnName] = time.Now()
	}

	if len(filter) == 0 || len(update) == 0 {
		return filter, update, errors.New("NilUpdateParams")
	}
	update = bson.M{"$set": update}

	return filter, update, nil
}
