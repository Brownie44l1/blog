package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB        *sqlx.DB
	JWTSecret string

	S3Bucket     string
	S3Region     string
	AWSAccessKey string
	AWSSecretKey string
}

func Load() *Config {
	// Try to load .env file - don't fail if it doesn't exist
	envPath := filepath.Join(".", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	} else {
		log.Println("✅ Loaded .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("❌ DATABASE_URL is not set. Please set it in .env file or export it")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalln("❌ JWT_SECRET is not set. Please set it in .env file or export it")
	}

	if len(jwtSecret) < 32 {
		log.Println("⚠️  WARNING: JWT_SECRET should be at least 32 characters for security")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln("❌ Failed to connect to DB:", err)
	}

	return &Config{
		DB:           db,
		JWTSecret:    jwtSecret,
		S3Bucket:     os.Getenv("S3_BUCKET"),
		S3Region:     os.Getenv("S3_REGION"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}
