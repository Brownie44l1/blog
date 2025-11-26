package repo

import (
	"strings"
	"github.com/jmoiron/sqlx"
    "github.com/Brownie44l1/blog/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *models.User) error {
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

func (r *UserRepo) GetUserByUsername(username string) (*models.User, error) {
	searchPattern := "%" + strings.ToLower(username) + "%"

	query := `
		SELECT * FROM blogs
		WHERE LOWER(title) ILIKE $1
		ORDER BY created_at DESC`

	// Pass the constructed pattern as the parameter for $1
	err := r.db.Select(&user, query, searchPattern)
	if err != nil {
		log.Printf("Error searching blogs for query %s: %v", searchQuery, err)
	}
	return &models.User{}, err
}

func (r *UserRepo) GetBlogCountByUserID(userID string) {
	
}
