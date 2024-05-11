package service

import (
	"blog-app/internal/models"
	"context"
	"time"
)

type IPostRepository interface {
	Save(context.Context, time.Duration, *models.Post) (*models.Post, error)
	Delete(context.Context, time.Duration, int) error
	GetAll(context.Context, time.Duration) ([]*models.Post, error)
	Get(context.Context, time.Duration, int) (*models.Post, error)
}

type PostService struct {
	repo     IPostRepository
	timeout  time.Duration
	telegram ITelegramService
}

func NewPostService(repo IPostRepository, telegram ITelegramService, timeout time.Duration) *PostService {
	return &PostService{repo: repo, timeout: timeout, telegram: telegram}
}

func (s *PostService) Create(ctx context.Context, post *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.repo.Save(ctx, s.timeout, post)
	if err != nil {
		return nil, err
	}
	err = s.telegram.Send(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
