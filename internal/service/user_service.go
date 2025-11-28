package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Brownie44l1/blog/internal/auth"
	"github.com/Brownie44l1/blog/internal/models"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already taken")
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetBlogCountByUserID(userID int64) (int, error)
}

type UserService interface {
	RegisterUser(username, password string) (*models.User, error)
	Authenticate(username, password string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	GetUserProfile(id int64) (*UserProfile, error)
}

type userService struct {
	userRepo UserRepository
}

type UserProfile struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	BlogCount int    `json:"blog_count"`
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(username, password string) (*models.User, error) {
	// Validate input
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	// Check if username exists
	_, err := s.userRepo.GetUserByUsername(username)
	if err == nil {
		return nil, ErrUsernameTaken
	}
	// Only proceed if error is "not found"
	if !errors.Is(err, sql.ErrNoRows) && err.Error() != fmt.Sprintf("user not found for username '%s'", username) {
		return nil, fmt.Errorf("error checking username: %w", err)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !auth.VerifyPassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

func (s *userService) GetUserProfile(id int64) (*UserProfile, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	blogCount, err := s.userRepo.GetBlogCountByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting blog count: %w", err)
	}

	return &UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		BlogCount: blogCount,
	}, nil
}