package main

import (
	"log"
	"net/http"

	"github.com/Brownie44l1/blog/config"
	"github.com/Brownie44l1/blog/internal/api"
	"github.com/Brownie44l1/blog/internal/repo"
	"github.com/Brownie44l1/blog/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()
	defer cfg.DB.Close()
	log.Println("✅ Connected to database!")

	// Initialize repositories
	userRepo := repo.NewUserRepo(cfg.DB)
	blogRepo := repo.NewBlogRepo(cfg.DB)
	log.Println("✅ Repositories initialized!")

	// Initialize services
	userService := service.NewUserService(userRepo)
	blogService := service.NewBlogService(blogRepo)
	log.Println("✅ Services initialized!")

	// Setup routes with all handlers
	router := api.SetupRoutes(userService, blogService, cfg.JWTSecret)
	log.Println("✅ Routes configured!")

	// Start server
	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}