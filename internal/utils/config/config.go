package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbHost     string
	//DbPort     string
	ServerPort string
	Secret     string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) InitENV() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.DbName = os.Getenv("DB_NAME")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	cfg.DbHost = os.Getenv("DB_HOST")
	//cfg.DbPort = os.Getenv("DB_PORT")
	cfg.Secret = os.Getenv("SECRET")
}
