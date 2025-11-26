package main

import (
	"log"
	"net/http"

	"myblog/internal/api"
	"myblog/internal/app"
	"myblog/internal/config"
	"myblog/internal/repo"
	"myblog/internal/service"
)

func main() {
    // Run migrations before connecting
    if err := app.RunMigrations(); err != nil {
        log.Fatalf("❌ Migration failed: %v", err)
    }
	log.Println("✅ Database migrations completed successfully!")

    db := config.NewDB()
    defer db.Close()
    log.Println("✅ Connected to database!")

    userRepo := repo.NewUserRepo(db)
    blogRepo := repo.NewBlogRepo(db)
    log.Println("✅ Queried Database Successfully!")

    userService := service.NewUserService(userRepo)
    blogService := service.NewBlogService(blogRepo, userRepo)
    log.Println("✅ Logic Applied Successfully!")

    userHandler := api.NewUserHandler(userService)
    blogHandler := api.NewBlogHandler(blogService)
    log.Println("✅ Parsed to and/or from JSON Successfully!")

    router := api.SetUpRoutes(userHandler, blogHandler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}