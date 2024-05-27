package buttons_repository

import (
	"NotSmokeBot/config"
	"NotSmokeBot/internal/model"
	"NotSmokeBot/internal/store/bson_queries"
	"NotSmokeBot/pkg/constant/status"
	"NotSmokeBot/pkg/dependences/tracer"
	"NotSmokeBot/pkg/errors/errlst"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			bson_queries.TgIdColumnName:         sentMessage.Sender,
			bson_queries.LastMessageColumnName:  sentMessage.Message,
			bson_queries.ProgressColumnName:     0,
			bson_queries.QuotationColumnName:    0,
			bson_queries.ActionStatusColumnName: status.StatusCasual,
			bson_queries.NickColumnName:         "",
			bson_queries.MoneyColumnName:        0.0,
			bson_queries.ChatIdColumnName:       sentMessage.ChatId,
			bson_queries.CreatedAtColumnName:    time.Now(),
			bson_queries.UpdatedAtColumnName:    time.Now(),
			bson_queries.DeletedAtColumnName:    nil,
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.AlreadyExistsError, err, "ButtonsMNGRepository.InsertNewUser.IsDuplicateKeyError")
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
		return tracer.SpanSetErrWrap(span, errlst.NilUpdateFieldsError, err, "ButtonsMNGRepository.UpdateUserByIds.getUpdateParams")
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
	filter[bson_queries.DeletedAtColumnName] = nil

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
	if updateUserInfo.ActionStatus != "" {
		update[bson_queries.ActionStatusColumnName] = updateUserInfo.ActionStatus
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

	if len(filter) == 1 || len(update) == 0 {
		return filter, update, errors.New("NilFields")
	}
	update = bson.M{"$set": update}

	return filter, update, nil
}

func (r *ButtonsMNGRepository) FindUsersByFilter(ctx context.Context, findUsersByFilter FindUsersByFilter) ([]model.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonsMNGRepository.FindUsersByFilter")
	defer span.End()

	collection := r.db.Collection(bson_queries.UsersCollection)
	filter, opts := getFindParams(findUsersByFilter)

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return []model.User{}, tracer.SpanSetErrWrap(span, errlst.ServerError, err, "ButtonsMNGRepository.FindUsersByFilter.Find")
	}
	defer cursor.Close(context.TODO())

	var users []User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return []model.User{}, tracer.SpanSetErrWrap(span, errlst.ServerError, err, "ButtonsMNGRepository.FindUsersByFilter.All")
	}

	if len(users) == 0 {
		return []model.User{}, tracer.SpanSetErrWrap(span, errlst.NothingFoundError, errors.New("NothingFoundError"), "ButtonsMNGRepository.FindUsersByFilter.len(users)")
	}

	var modelUsers []model.User
	for _, user := range users {
		modelUsers = append(modelUsers, user.toUserModel())
	}

	return modelUsers, nil
}

func getFindParams(findUsersByFilter FindUsersByFilter) (bson.M, *options.FindOptions) {
	filter := bson.M{}

	if len(findUsersByFilter.MongoId) != 0 {
		filter[bson_queries.IdColumnName] = bson.M{"$in": findUsersByFilter.MongoId}
	}
	if len(findUsersByFilter.TgId) != 0 {
		filter[bson_queries.TgIdColumnName] = bson.M{"$in": findUsersByFilter.TgId}
	}
	if len(findUsersByFilter.LastMessage) != 0 {
		filter[bson_queries.LastMessageColumnName] = bson.M{"$in": findUsersByFilter.LastMessage}
	}

	progressFilter := bson.M{}
	if len(findUsersByFilter.Progress) != 0 {
		progressFilter["$in"] = findUsersByFilter.Progress
	}
	if findUsersByFilter.ProgressLess != nil {
		progressFilter["$lte"] = findUsersByFilter.ProgressLess
	}
	if findUsersByFilter.ProgressMore != nil {
		progressFilter["$gte"] = findUsersByFilter.ProgressMore
	}
	if len(progressFilter) != 0 {
		filter[bson_queries.ProgressColumnName] = progressFilter
	}

	quotationFilter := bson.M{}
	if len(findUsersByFilter.Quotation) != 0 {
		quotationFilter["$in"] = findUsersByFilter.Quotation
	}
	if findUsersByFilter.QuotationLess != nil {
		quotationFilter["$lte"] = findUsersByFilter.QuotationLess
	}
	if findUsersByFilter.QuotationMore != nil {
		quotationFilter["$gte"] = findUsersByFilter.QuotationMore
	}
	if len(quotationFilter) != 0 {
		filter[bson_queries.QuotationColumnName] = quotationFilter
	}

	if len(findUsersByFilter.ActionStatus) != 0 {
		filter[bson_queries.ActionStatusColumnName] = bson.M{"$in": findUsersByFilter.ActionStatus}
	}
	if len(findUsersByFilter.Nick) != 0 {
		filter[bson_queries.NickColumnName] = bson.M{"$in": findUsersByFilter.Nick}
	}

	moneyFilter := bson.M{}
	if len(findUsersByFilter.Money) != 0 {
		moneyFilter["$in"] = findUsersByFilter.Money
	}
	if findUsersByFilter.MoneyLess != nil {
		moneyFilter["$lte"] = findUsersByFilter.MoneyLess
	}
	if findUsersByFilter.MoneyMore != nil {
		moneyFilter["$gte"] = findUsersByFilter.MoneyMore
	}
	if len(moneyFilter) != 0 {
		filter[bson_queries.MoneyColumnName] = moneyFilter
	}

	if len(findUsersByFilter.ChatId) != 0 {
		filter[bson_queries.ChatIdColumnName] = bson.M{"$in": findUsersByFilter.ChatId}
	}

	createdAtFilter := bson.M{}
	if !findUsersByFilter.CreatedAtLess.IsZero() {
		createdAtFilter["$lte"] = findUsersByFilter.CreatedAtLess
	}
	if !findUsersByFilter.CreatedAtMore.IsZero() {
		createdAtFilter["$gte"] = findUsersByFilter.CreatedAtMore
	}
	if len(createdAtFilter) != 0 {
		filter[bson_queries.CreatedAtColumnName] = createdAtFilter
	}

	updatedAtFilter := bson.M{}
	if !findUsersByFilter.UpdatedAtLess.IsZero() {
		updatedAtFilter["$lte"] = findUsersByFilter.UpdatedAtLess
	}
	if !findUsersByFilter.UpdatedAtMore.IsZero() {
		updatedAtFilter["$gte"] = findUsersByFilter.UpdatedAtMore
	}
	if len(updatedAtFilter) != 0 {
		filter[bson_queries.UpdatedAtColumnName] = updatedAtFilter
	}

	deletedAtFilter := bson.M{}
	if findUsersByFilter.DeletedAt != false {
		deletedAtFilter["$ne"] = nil
	}
	if !findUsersByFilter.DeletedAtLess.IsZero() {
		deletedAtFilter["$lte"] = findUsersByFilter.DeletedAtLess
	}
	if !findUsersByFilter.DeletedAtMore.IsZero() {
		deletedAtFilter["$gte"] = findUsersByFilter.DeletedAtMore
	}
	if len(deletedAtFilter) != 0 {
		filter[bson_queries.DeletedAtColumnName] = deletedAtFilter
	} else {
		filter[bson_queries.DeletedAtColumnName] = nil
	}

	findOptions := options.Find()
	if findUsersByFilter.Offset != 0 {
		findOptions.SetSkip(findUsersByFilter.Offset)
	}
	if findUsersByFilter.Limit != 0 {
		findOptions.SetLimit(findUsersByFilter.Limit)
	}

	return filter, findOptions
}
