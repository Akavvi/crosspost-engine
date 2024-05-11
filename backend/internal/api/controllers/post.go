package controllers

import (
	"blog-app/internal/models"
	"blog-app/pkg/service"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type PostService interface {
	Create(post *models.Post) (*models.Post, error)
	Delete(id int) error
	Find(id int) (*models.Post, error)
	GetAll() ([]*models.Post, error)
}

type PostController struct {
	service service.PostService
	timeout time.Duration
}

func NewPostController(service *service.PostService) *PostController {
	return &PostController{
		service: *service,
	}
}

func (c *PostController) AddPost(w http.ResponseWriter, r *http.Request) {
	p := &models.Post{}
	err := json.NewDecoder(r.Body).Decode(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.BeforeCreate()

	post, err := c.service.Create(r.Context(), p)
	if err != nil {
		log.Printf("%v", err)
		json.NewEncoder(w).Encode(map[string]bool{"created": false})
		return
	}
	json.NewEncoder(w).Encode(post)
}
