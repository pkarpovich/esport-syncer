package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	PandaScore struct {
		TeamId string `env:"PANDASCORE_TEAM_ID" env-default:"6224"`
		ApiKey string `env:"PANDASCORE_API_KEY"`
	}
	Calendar struct {
		Name            string `env:"CALENDAR_NAME" env-default:"Esport matches"`
		Color           string `env:"CALENDAR_COLOR" env-default:"red"`
		RefreshInterval string `env:"CALENDAR_REFRESH_INTERVAL" env-default:"P1D"`
	}
	Port       string `env:"PORT" env-default:"1710"`
	Secret     string `env:"SECRET_KEY"`
	ConfigPath string `env:"CONFIG_PATH" env-default:"./.config/config.json"`
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[WARN] error while loading .env file: %v", err)
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
