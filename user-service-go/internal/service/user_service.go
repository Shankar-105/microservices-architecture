package service

import (
	"context"
	"errors"

	"github.com/shank/bookstore-microservices/user-service-go/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo *repository.PostgresUserRepo
}

func NewUserService(repo *repository.PostgresUserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*repository.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
