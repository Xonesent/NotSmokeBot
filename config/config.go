package config

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
)

const configPath = "./config/config.json"

type ConfigPg struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	DbName   string `validate:"required"`
	Password string `validate:"required"`
	SSLMode  string `validate:"required"`
}

type Config struct {
	Postgres ConfigPg
}

func LoadConfig() (c *Config, err error) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(jsonFile).Decode(&c)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}
	return
}
