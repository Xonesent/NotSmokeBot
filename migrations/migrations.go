package main

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/tools/logger"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatalf("Error to init logger: %v\n", err)
	}

	app := &cli.App{
		Name:  "migrate",
		Usage: "Database migration tool",
		Commands: []*cli.Command{
			{
				Name:   "up",
				Usage:  "Apply all up migrations",
				Action: migrateUp,
			},
			{
				Name:   "down",
				Usage:  "Apply all down migrations",
				Action: migrateDown,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 || c.NArg() > 1 {
				cli.ShowAppHelp(c)
				return cli.Exit("Invalid number of arguments or unrecognized flag", 1)
			}

			arg := c.Args().Get(0)
			if arg != "up" && arg != "down" {
				cli.ShowAppHelp(c)
				return cli.Exit("Invalid flag: "+arg, 1)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		zap.L().Error("Impossible to migrate database", zap.Error(err))
	}
}

func migrateUp(c *cli.Context) error {
	zap.L().Info("Applying up migrations...")

	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	zap.L().Info("Migration up completed successfully!")
	return nil
}

func migrateDown(c *cli.Context) error {
	zap.L().Info("Applying down migrations...")

	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	zap.L().Info("Migration down completed successfully!")
	return nil
}

func getMigrateInstance() (*migrate.Migrate, error) {
	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("Error loading env variables", zap.Error(err))
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Error loading config", zap.Error(err))
	}

	credential := options.Credential{
		Username: cfg.Mongo.User,
		Password: cfg.Mongo.Password,
	}
	connUrl := fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)

	clientOptions := options.Client().ApplyURI(connUrl).SetAuth(credential)
	client, err := mongo.Connect(context.Background(), clientOptions)
	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: cfg.Mongo.Database})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.Mongo.Database,
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
