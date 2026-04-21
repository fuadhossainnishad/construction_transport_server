package config

import (
	"construction_transport_server/pkg/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db  DBConfig
	App AppConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		App: AppConfig{
			App_port: utils.StringToInt(os.Getenv("PORT"), 3000),
		},
		Db: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     utils.StringToInt(os.Getenv("DB_PORT"), 5432),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}
}
