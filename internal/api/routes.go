package api

import (
    "net/http"
    
    "github.com/Brownie44l1/blog/internal/middleware"
    "github.com/Brownie44l1/blog/internal/services"
)

func SetupRoutes(
    userService services.UserService,
    blogService services.BlogService,
    jwtSecret string,
) *http.ServeMux {
    mux := http.NewServeMux()
    
    // Create handlers
    authHandler := NewAuthHandler(userService, jwtSecret)
    blogHandler := NewBlogHandler(blogService)
    
    // Create middleware
    authMiddleware := middleware.AuthMiddleware(jwtSecret)
    
    // Public routes (no auth needed)
    mux.HandleFunc("/register", authHandler.Register)
    mux.HandleFunc("/login", authHandler.Login)
    // How would you add GET /blogs?
    
    // Protected routes (need auth)
    // This is the tricky part - you need to wrap the handler with middleware
    mux.Handle("/blogs/create", authMiddleware(http.HandlerFunc(blogHandler.CreateBlog)))
    
    return mux
}