package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const path = "./config"

type Config struct {
	OpenTelemetry struct {
		URL         string `validate:"required"`
		ServiceName string `validate:"required"`
	}
	Mongo struct {
		Host     string `validate:"required"`
		Port     string `validate:"required"`
		User     string `validate:"required"`
		Password string `validate:"required"`
	}
	Telegram struct {
		Token string `validate:"required"`
	}
}

func LoadConfig() (c *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err = viper.Unmarshal(&c); err != nil {
		return c, err
	}

	if err = validator.New().Struct(c); err != nil {
		return c, err
	}
	return
}
