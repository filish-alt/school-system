package config

import (
	"os"
)

type Config struct {
	DBPath    string
	JWTSecret string
	Port      string
}

func Load() Config {
	db := os.Getenv("DB_PATH")
	if db == "" {
		db = "school.db"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		DBPath:    db,
		JWTSecret: secret,
		Port:      port,
	}
}

