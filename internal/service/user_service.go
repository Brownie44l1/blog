package service

import (
	"github.com/Brownie44l1/blog/internal/models"
	"github.com/Brownie44l1/blog/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *UserService) Register(username, password string) (*models.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}
	err1 := s.UserRepo.Create(user)
	if err1 != nil {
		return nil, err1
	}
	return user, nil
}

func (s *UserService) verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil 
}

func (s *UserService) GetProfile(id string) (*models.User, error) {
	return  s.UserRepo.GetByID(id)
}
