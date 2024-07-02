package service

import (
	"blog-app/internal/models"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PostRepository interface {
	Save(context.Context, time.Duration, *models.Post) (*models.Post, error)
	Update(context.Context, time.Duration, *models.Post) (*models.Post, error)
	Delete(context.Context, time.Duration, int) error
	GetAll(context.Context, time.Duration) ([]*models.Post, error)
	Get(context.Context, time.Duration, int) (*models.Post, error)
}

type ITelegramService interface {
	Send(post *models.Post) error
}

type PostService struct {
	repo     PostRepository
	timeout  time.Duration
	telegram TelegramService
}

func NewPostService(repo PostRepository, timeout time.Duration, telegram TelegramService) *PostService {
	return &PostService{repo: repo, timeout: timeout, telegram: telegram}
}

func (s *PostService) Create(ctx context.Context, r *http.Request) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post := &models.Post{}
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")

	file, header, _ := r.FormFile("file")
	if file != nil {
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		var buf bytes.Buffer
		_, err := io.Copy(&buf, file)
		if err != nil {
			return nil, err
		}

		name := uuid.NewString()
		extension := filepath.Ext(header.Filename)
		if extension == "" {
			extension = ".png"
		}
		filename := fmt.Sprintf("%s%s", name, extension)
		err = os.WriteFile(fmt.Sprintf("attachments/%s", filename), buf.Bytes(), 0644)
		buf.Reset()
		if err != nil {
			log.Print(err)
		}
		post.Attachment = &filename
	}
	if s.telegram.Alive == true {
		message, err := s.telegram.Send(post)
		if err != nil {
			log.Println("Telegram is not responding or service is not initialized")
		}
		post.TelegramMessageId = &message.MessageId
	}
	post.BeforeCreate()
	p, err := s.repo.Save(ctx, s.timeout, post)
	if err != nil {
		if s.telegram.Alive == true {
			_ = s.telegram.Delete(*post.TelegramMessageId)
		}
		return nil, err
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
	err = s.telegram.Delete(*exists.TelegramMessageId)
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

func (s *PostService) Update(ctx context.Context, id int, r *http.Request) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post, err := s.repo.Get(ctx, s.timeout, id)
	if err != nil || post == nil {
		return nil, err
	}

	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	post.UpdatedAt = time.Now()

	p, err := s.repo.Update(ctx, s.timeout, post)
	if err != nil {
		return nil, err
	}

	if p.TelegramMessageId != nil {
		err = s.telegram.Update(*p.TelegramMessageId, p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}
