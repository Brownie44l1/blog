package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
    DB        *sqlx.DB
    JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load();

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatalln("JWT_SECRET is not set")
    }

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}
	return &Config{
        DB:        db,
        JWTSecret: jwtSecret,
    }
}
