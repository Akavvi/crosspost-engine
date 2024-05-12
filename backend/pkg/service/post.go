package service

import (
	"blog-app/internal/models"
	"context"
	"log"
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

func NewPostService(repo IPostRepository, timeout time.Duration, telegram ITelegramService) *PostService {
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
		log.Println("Telegram is not responding or service is not initialized")
	}
	return p, nil
}

func (s *PostService) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	exists, err := s.repo.Get(ctx, s.timeout, id)
	if err != nil || exists == nil {
		return err
	}
	err = s.repo.Delete(ctx, s.timeout, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) Find(ctx context.Context, id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	post, err := s.repo.Get(ctx, s.timeout, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAll(ctx context.Context) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	posts, err := s.repo.GetAll(ctx, s.timeout)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
