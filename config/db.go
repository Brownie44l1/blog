package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewDB() *sqlx.DB {
	if err := godotenv.Load(); err != nil {
		// If not found in current directory, try parent directory
		if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

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
