package route

import (
	"blog-app/internal/api/controllers"
	"blog-app/internal/repository"
	"blog-app/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewPostRoute(db *sqlx.DB, router *mux.Router, telegram service.ITelegramService) {
	repo := repository.NewPostRepository(db)
	s := service.NewPostService(repo, telegram, 5*time.Second)
	c := controllers.NewPostController(s)

	router.HandleFunc("/addPost", c.AddPost).Methods("POST")
}
