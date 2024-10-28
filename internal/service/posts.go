package service

import (
	"context"
	"newsSite/internal/model"
)

func (s Service) GetAllPosts(ctx context.Context) (model.PostsList, error) {
	return s.repo.GetAllPosts(ctx)
}

func (s Service) AddPost(ctx context.Context, post model.Posts) error {
	return s.repo.AddPost(ctx, post)
}
