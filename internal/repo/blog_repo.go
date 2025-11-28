package repo

import (
	"fmt" 
	"log" 
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/Brownie44l1/blog/internal/models"
)

type BlogRepo struct {
	db *sqlx.DB
}

func NewBlogRepo(db *sqlx.DB) *BlogRepo {
	return &BlogRepo{db: db}
}

func (r *BlogRepo) CreateBlog(blog *models.Blog) error {
	query := `
		INSERT INTO blogs (user_id, title, content)
		VALUES($1, $2, $3)
		RETURNING id, created_at`
	return r.db.QueryRow(
		query, blog.UserId, blog.Title, blog.Content,
	).Scan(&blog.ID, &blog.CreatedAt)
}

func (r *BlogRepo) GetBlogByID(id string) (*models.Blog, error) {
	var blog models.Blog
	query := `SELECT * FROM blogs WHERE id=$1`
	err := r.db.Get(&blog, query, id)
	if err != nil {
		log.Printf("Error getting blog by ID %s: %v", id, err)
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepo) GetBlogByUserID(userID string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT * FROM blogs WHERE user_id=$1 ORDER BY created_at DESC`
	err := r.db.Select(&blogs, query, userID)
	if err != nil {
		log.Printf("Error getting blogs for user %s: %v", userID, err)
	}
	return blogs, err
}

func (r *BlogRepo) DeleteBlog(blogID, userID string) error {
	query := `
		DELETE FROM blogs
		WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, blogID, userID)
	if err != nil {
		log.Printf("Error deleting blog %s by user %s: %v", blogID, userID, err)
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows after delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no blog found with ID %s for user %s, or user is not the owner", blogID, userID)
	}

	return nil
}

func (r *BlogRepo) GetAllBlogs(limit, offset string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `
		SELECT * FROM blogs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	err := r.db.Select(&blogs, query, limit, offset)
	if err != nil {
		log.Printf("Error getting all blogs with limit %s offset %s: %v", limit, offset, err)
	}
	return blogs, err
}

func (r *BlogRepo) SearchBlogs(searchQuery string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	searchPattern := "%" + strings.ToLower(searchQuery) + "%"
	query := `
		SELECT * FROM blogs
		WHERE LOWER(title) ILIKE $1 OR LOWER(content) ILIKE $1
		ORDER BY created_at DESC`

	err := r.db.Select(&blogs, query, searchPattern)
	if err != nil {
		log.Printf("Error searching blogs for query %s: %v", searchQuery, err)
	}
	return blogs, err
}

