package service

import (
    "database/sql"
    "errors"
    "fmt"

    "github.com/Brownie44l1/blog/internal/auth"
    "github.com/Brownie44l1/blog/internal/models"
    "github.com/Brownie44l1/blog/internal/repo"
)

var (
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUsernameTaken      = errors.New("username already taken")
)

type UserService interface {
    RegisterUser(username, password string) (*models.User, error)
    Authenticate(username, password string) (*models.User, error)
    GetUserByID(id int64) (*models.User, error)
}

type userService struct {
    userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *userService {
    return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(username, password string) (*models.User, error) {
    // Check if username exists
    _, err := s.userRepo.GetUserByUsername(username)
    if err == nil {
        return nil, ErrUsernameTaken
    }
    // Only proceed if error is "not found"
    if !errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("error checking username: %w", err)
    }

    // Use YOUR auth.HashPassword function here
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

    user.Password = "" // Clear before returning
    return user, nil
}

func (s *userService) Authenticate(username, password string) (*models.User, error) {
    user, err := s.userRepo.GetUserByUsername(username)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrInvalidCredentials
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    // Use YOUR auth.VerifyPassword function here
    if !auth.VerifyPassword(user.Password, password) {
        return nil, ErrInvalidCredentials
    }

    user.Password = ""
    return user, nil
}