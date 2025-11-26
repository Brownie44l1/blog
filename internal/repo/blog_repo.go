package repo

import (
	"myblog/internal/models"
	"github.com/jmoiron/sqlx"
)

type BlogRepo struct {
	db *sqlx.DB
}

func NewBlogRepo(db *sqlx.DB) *BlogRepo {
	return &BlogRepo{db: db}
}

func (r *BlogRepo) Create(blog *models.Blog) error {
	query := `
		INSERT INTO blogs (user_id, title, content)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(
		query, blog.UserId, blog.Title, blog.Content,
	).Scan(&blog.ID, &blog.CreatedAt, &blog.UpdatedAt)
}

func (r *BlogRepo) GetByID(id string) (*models.Blog, error) {
	var blog models.Blog
	query := `SELECT * FROM blogs WHERE id=$1`
	err := r.db.Get(&blog, query, id)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepo) ListByUser(userID string) ([]models.Blog, error) {
	blogs := []models.Blog{}
	query := `SELECT * FROM blogs WHERE user_id=$1 ORDER BY created_at DESC`
	err := r.db.Select(&blogs, query, userID) 
	return blogs, err
}
