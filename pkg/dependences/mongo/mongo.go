package mongodb

import (
	"NotSmokeBot/config"
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg *config.Config) (*mongo.Client, error) {
	credential := options.Credential{
		Username: cfg.Mongo.User,
		Password: cfg.Mongo.Password,
	}
	connUrl := fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)
	clientOptions := options.Client().ApplyURI(connUrl).SetAuth(credential)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
