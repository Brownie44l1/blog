package api

import (
	"net/http"

	"github.com/Brownie44l1/blog/internal/middleware"
	"github.com/Brownie44l1/blog/internal/service"
)

func SetupRoutes(
	userService service.UserService,
	blogService service.BlogService,
	jwtSecret string,
) *http.ServeMux {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(userService, jwtSecret)
	blogHandler := NewBlogHandler(blogService)
	userHandler := NewUserHandler(userService)

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// ==================== AUTH ROUTES ====================
	// Public routes - no authentication required
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// ==================== USER ROUTES ====================
	// Get authenticated user's profile (protected)
	mux.Handle("/users/me", authMiddleware(http.HandlerFunc(userHandler.GetMe)))

	// Get any user's profile (public)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a user blog request: /users/{id}/blogs
		if len(r.URL.Path) > 7 && r.URL.Path[len(r.URL.Path)-6:] == "/blogs" {
			blogHandler.GetUserBlogs(w, r)
			return
		}
		// Otherwise, it's a user profile request: /users/{id}
		userHandler.GetProfile(w, r)
	})

	// ==================== BLOG ROUTES ====================
	// Create blog (protected)
	mux.Handle("/blogs/create", authMiddleware(http.HandlerFunc(blogHandler.CreateBlog)))

	// Get authenticated user's blogs (protected)
	mux.Handle("/blogs/me", authMiddleware(http.HandlerFunc(blogHandler.GetMyBlogs)))

	// Search blogs (public)
	mux.HandleFunc("/blogs/search", blogHandler.SearchBlogs)

	// Blog operations by ID
	mux.HandleFunc("/blogs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Public: anyone can view a blog
			blogHandler.GetBlog(w, r)
		case http.MethodDelete:
			// Protected: only owner can delete
			authMiddleware(http.HandlerFunc(blogHandler.DeleteBlog)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// List all blogs with pagination (public)
	mux.HandleFunc("/blogs", blogHandler.ListBlogs)

	return mux
}