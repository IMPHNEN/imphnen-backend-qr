package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Fallback to .env file, ignore error if not found
	_ = viper.ReadInConfig()

	cfg := &Config{
		Port:               viper.GetString("PORT"),
		DatabaseURL:        viper.GetString("DATABASE_URL"),
		JWTSecret:          viper.GetString("JWT_SECRET"),
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}
