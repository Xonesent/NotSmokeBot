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

func (r *ButtonsMNGRepository) InsertNewUser(ctx context.Context, startMessage StartMessage) (primitive.ObjectID, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ButtonsMNGRepository.InsertNewUser")
	defer span.End()

	collection := r.db.Collection(bson_queries.UsersCollection)
	insertResult, err := collection.InsertOne(context.TODO(),
		bson.M{
			bson_queries.TdIdColumnName:        startMessage.Sender,
			bson_queries.LastMessageColumnName: startMessage.Message,
			bson_queries.ProgressColumnName:    0,
			bson_queries.QuotationColumnName:   0,
			bson_queries.NickColumnName:        "",
			bson_queries.MoneyColumnName:       0.0,
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.AlreadyExists, err, "ButtonsMNGRepository.InsertNewUser.InsertOne")
		}
		return primitive.ObjectID{}, err
	}

	value, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, tracer.SpanSetErrWrap(span, errlst.ConvertionError, errors.New("primitive.ObjectID convert error"), "ButtonsMNGRepository.InsertNewUser.ok")
	}

	return value, nil
}
