package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewDB() *sqlx.DB {
	_ = godotenv.Load();

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("DATABASE_URL is not set")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}
	return db
}
