package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
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
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func migrateUp(c *cli.Context) error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migration up completed successfully!")
	return nil
}

func migrateDown(c *cli.Context) error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migration down completed successfully!")
	return nil
}

func getMigrateInstance() (*migrate.Migrate, error) {
	mongoURI := "mongodb://localhost:27017"
	dbName := "mydatabase5"

	credential := options.Credential{
		Username: "TelegramBot",
		Password: "c2149ca6-a9e6-49d3-9a65-22e48e7ae461",
	}

	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(credential)

	client, err := mongo.Connect(context.Background(), clientOptions)

	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
