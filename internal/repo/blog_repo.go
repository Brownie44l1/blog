package repo

import (
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
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepo) GetBlogByUserID(userID string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT * FROM blogs WHERE user_id=$1 ORDER BY created_at DESC`
	err := r.db.Select(&blogs, query, userID) 
	return blogs, err
}

func (r *BlogRepo) DeleteBlog(blogID, userID string) {
	query := `SELECT * FROM blogs WHERE id=$1`
	
}

func (r *BlogRepo) GetAllBlogs(limit, offset string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT * FROM blogs ORDER BY created_at DESC LIMIT 5 OFFSET `
	err := r.db.Select(&blogs, query, limit, offset) 
	return blogs, err
}

func (r *BlogRepo) SearchBlogs(query string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT * FROM blogs WHERE title LIKE="%query%" ORDER BY created_at DESC`
	err := r.db.Select(&blogs, query, query) 
	return blogs, err
}
