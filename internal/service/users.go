package service

import (
	"context"
	"newsSite/internal/model"
)

func (s Service) RegisterUser(ctx context.Context, user model.User) error {
	return s.repo.RegisterUser(ctx, user)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
