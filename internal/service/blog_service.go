package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/Brownie44l1/blog/internal/models"
)

// BlogRepository defines the interface for blog data operations.
type BlogRepository interface {
	CreateBlog(blog *models.Blog) error
	GetBlogByID(id int64) (*models.Blog, error)
	GetBlogByUserID(userID int64) ([]models.Blog, error)
	DeleteBlog(blogID, userID int64) error
	UpdateBlog(blog *models.Blog) error 
	GetAllBlogs(limit, offset int64) ([]models.Blog, error)
	SearchBlogs(searchQuery string) ([]models.Blog, error)
}

// BlogService defines the interface for blog business logic
type BlogService interface {
	Create(blog *models.Blog) error
	GetByID(id int64) (*models.Blog, error)
	GetByUserID(userID int64) ([]models.Blog, error)
	Update(blog *models.Blog) error
	Delete(blogID, userID int64) error
	ListAll(limit, offset int64) ([]models.Blog, error)
	Search(query string) ([]models.Blog, error)
}

// blogService is the concrete implementation
type blogService struct {
	repo BlogRepository
}

// NewBlogService creates a new BlogService instance.
func NewBlogService(r BlogRepository) BlogService {
	return &blogService{repo: r}
}

// Create validates and creates a new blog post.
func (s *blogService) Create(blog *models.Blog) error {
	if strings.TrimSpace(blog.Title) == "" {
		return fmt.Errorf("blog title cannot be empty")
	}
	if strings.TrimSpace(blog.Content) == "" {
		return fmt.Errorf("blog content cannot be empty")
	}

	if err := s.repo.CreateBlog(blog); err != nil {
		log.Printf("Service error creating blog: %v", err)
		return fmt.Errorf("failed to create blog post: %w", err)
	}
	return nil
}

// GetByID retrieves a single blog post by ID.
func (s *blogService) GetByID(id int64) (*models.Blog, error) {
	blog, err := s.repo.GetBlogByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving blog ID %d: %w", id, err)
	}
	return blog, nil
}

// GetByUserID retrieves all blogs for a specific user.
func (s *blogService) GetByUserID(userID int64) ([]models.Blog, error) {
	blogs, err := s.repo.GetBlogByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving blogs for user %d: %w", userID, err)
	}
	return blogs, nil
}

func (s *blogService) Update(blog *models.Blog) error {
	err := s.repo.UpdateBlog(blog)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// Delete handles blog deletion with ownership verification.
func (s *blogService) Delete(blogID, userID int64) error {
	err := s.repo.DeleteBlog(blogID, userID)
	if err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}
	return nil
}

// ListAll retrieves all blogs with pagination.
func (s *blogService) ListAll(limit, offset int64) ([]models.Blog, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	blogs, err := s.repo.GetAllBlogs(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing all blogs: %w", err)
	}
	return blogs, nil
}

// Search queries blogs based on a search term.
func (s *blogService) Search(query string) ([]models.Blog, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	blogs, err := s.repo.SearchBlogs(query)
	if err != nil {
		return nil, fmt.Errorf("error during blog search: %w", err)
	}
	return blogs, nil
}