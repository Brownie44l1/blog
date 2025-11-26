package app

import (
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/joho/godotenv"
)

func RunMigrations() error {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ No .env file found, using system env instead")
    }

    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return fmt.Errorf("DATABASE_URL not set")
    }

    m, err := migrate.New(
        "file://C:/Users/HP/GoLang/myblog/migrations",
        dbURL,
    )
    if err != nil {
        return err
    }
    defer m.Close()
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}