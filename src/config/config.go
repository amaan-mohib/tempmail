package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AuthenticationConfig struct {
	JWTSecret string
}

type TConfig struct {
	Authentication AuthenticationConfig
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
	}
	Config = config
	return config
}
