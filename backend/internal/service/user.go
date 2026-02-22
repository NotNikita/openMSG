package service

import (
	"context"
	"errors"

	"app/internal/models"
	"app/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, nickname, publicKey, avatar string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, nickname, publicKey, avatar string) (*models.User, error) {
	if nickname == "" {
		return nil, errors.New("nickname is required")
	}
	if publicKey == "" {
		return nil, errors.New("public_key is required")
	}
	if avatar == "" {
		return nil, errors.New("avatar is required")
	}
	return s.repo.Create(ctx, nickname, publicKey, avatar)
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}
