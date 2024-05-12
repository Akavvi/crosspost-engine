package controllers

import (
	"blog-app/internal/models"
	"blog-app/pkg/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type PostService interface {
	Create(ctx context.Context, post *models.Post) (*models.Post, error)
	Delete(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.Post, error)
	GetAll(ctx context.Context) ([]*models.Post, error)
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

func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
	p := &models.Post{}
	err := json.NewDecoder(r.Body).Decode(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.BeforeCreate()

	post, err := c.service.Create(r.Context(), p)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]bool{"created": false})
		return
	}
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.service.Delete(r.Context(), id)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"deleted": "false", "err": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"deleted": "true"})
}

func (c *PostController) Find(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := c.service.Find(r.Context(), id)
	if err != nil || post == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := c.service.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
