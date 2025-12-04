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

func (r *BlogRepo) GetBlogByID(id int64) (*models.Blog, error) {
	var blog models.Blog
	query := `SELECT id, user_id, title, content, created_at FROM blogs WHERE id=$1`
	err := r.db.Get(&blog, query, id)
	if err != nil {
		log.Printf("Error getting blog by ID %d: %v", id, err)
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepo) GetBlogByUserID(userID int64) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT id, user_id, title, content, created_at FROM blogs WHERE user_id=$1 ORDER BY created_at DESC`
	err := r.db.Select(&blogs, query, userID)
	if err != nil {
		log.Printf("Error getting blogs for user %d: %v", userID, err)
	}
	return blogs, err
}

func (r *BlogRepo) UpdateBlog(blog *models.Blog) error {
    query := `
        UPDATE blogs 
        SET title = $1, content = $2, updated_at = NOW()
        WHERE id = $3 AND user_id = $4
        RETURNING id, user_id, title, content, created_at, updated_at
    `
    return r.db.QueryRow(query, blog.Title, blog.Content, blog.ID, blog.UserId).
        Scan(&blog.ID, &blog.UserId, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
}

func (r *BlogRepo) DeleteBlog(blogID, userID int64) error {
	query := `
		DELETE FROM blogs
		WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, blogID, userID)
	if err != nil {
		log.Printf("Error deleting blog %d by user %d: %v", blogID, userID, err)
		return fmt.Errorf("failed to delete blog: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows after delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no blog found with ID %d for user %d, or user is not the owner", blogID, userID)
	}

	return nil
}

func (r *BlogRepo) GetAllBlogs(limit, offset int64) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `
		SELECT id, user_id, title, content, created_at FROM blogs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	err := r.db.Select(&blogs, query, limit, offset)
	if err != nil {
		log.Printf("Error getting all blogs with limit %d offset %d: %v", limit, offset, err)
	}
	return blogs, err
}

func (r *BlogRepo) SearchBlogs(searchQuery string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	searchPattern := "%" + strings.ToLower(searchQuery) + "%"
	query := `
		SELECT id, user_id, title, content, created_at 
		FROM blogs
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY created_at DESC`

	err := r.db.Select(&blogs, query, searchPattern)
	if err != nil {
		log.Printf("Error searching blogs for query %s: %v", searchQuery, err)
	}
	return blogs, err
}

