package service

import (
	"errors"
    "github.com/Brownie44l1/blog/internal/models"
    "github.com/Brownie44l1/blog/internal/repo"
)

type BlogService struct {
	blogRepo *repo.BlogRepo
	userRepo *repo.UserRepo
}

func NewBlogService(blogRepo *repo.BlogRepo, userRepo *repo.UserRepo) *BlogService {
	return &BlogService{blogRepo: blogRepo, userRepo: userRepo,}
}

func (s *BlogService) Publish(userId, title, content string) (*models.Blog, error) {
	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, errors.New("User not found")
	}
	blog := &models.Blog{
		UserId: user.ID,
		Title: title,
		Content: content,
	}

	err1 := s.blogRepo.Create(blog)
	if err1 != nil {
		return nil, err
	}
	return blog, nil
}

func (s *BlogService) GetBlog(id string) (*models.Blog, error) {
	return s.blogRepo.GetByID(id)
}

func (s *BlogService) ListUserBlog(userID string) ([]models.Blog, error) {
	return s.blogRepo.ListByUser(userID)
}