package repo

import (
	"database/sql" 
	"fmt" 
	"log" 
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

func (r *UserRepo) GetByID(id int64) (*models.User, error) {
    query := `
        SELECT u.id, u.username, u.password, COALESCE(COUNT(b.id), 0) as blog_count
        FROM users u
        LEFT JOIN blogs b ON u.id = b.user_id
        WHERE u.id = $1
        GROUP BY u.id
    `
    user := &models.User{}
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.BlogCount)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	query := `
		SELECT * FROM users
		WHERE LOWER(username) = LOWER($1)`

	err := r.db.Get(&user, query, username)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found for username '%s'", username)
		}

		log.Printf("Database error getting user by username %s: %v", username, err)
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetBlogCountByUserID(userID int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(id) FROM blogs
		WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)

	if err != nil {
		log.Printf("Error counting blogs for user %d: %v", userID, err)
		return 0, fmt.Errorf("failed to count blogs for user %d: %w", userID, err)
	}
	
	return count, nil
}
