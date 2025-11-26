package repo

import (
	"github.com/jmoiron/sqlx"
    "github.com/Brownie44l1/blog/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES($1, $2)
		RETURNING id`
	return r.db.QueryRow(
		query, user.Username, user.Password,
	).Scan(&user.ID)
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
