package main

import (
	"log"
	"net/http"

	"github.com/Brownie44l1/blog/internal/api"
	"github.com/Brownie44l1/blog/config"
	"github.com/Brownie44l1/blog/internal/repo"
	"github.com/Brownie44l1/blog/internal/service"
)

func main() {
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