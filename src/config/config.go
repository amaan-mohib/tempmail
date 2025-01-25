package config

import (
	"log"
	"os"
	"tempgalias/src/utils"

	"github.com/joho/godotenv"
)

type AuthenticationConfig struct {
	JWTSecret string
}

type TConfig struct {
	Authentication AuthenticationConfig
	BaseEmail      string
	DefaultExpiry  int64
	DatabaseURL    string
}

var Config *TConfig

func LoadConfig() *TConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &TConfig{
		Authentication: AuthenticationConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		BaseEmail:     os.Getenv("BASE_EMAIL"),
		DefaultExpiry: utils.ParseInt(os.Getenv("DEFAULT_EXPIRY")),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
	}
	Config = config
	return config
}
