package repo

import (
	"github.com/jmoiron/sqlx"
    "myblog/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at, blog_count`
	return r.db.QueryRow(
		query, user.Username, user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.BlogCount)
}

func (r *UserRepo) GetByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id=$1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
