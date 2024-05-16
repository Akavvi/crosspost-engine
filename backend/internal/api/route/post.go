package route

import (
	"blog-app/internal/api/controllers"
	"blog-app/internal/repository"
	"blog-app/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewPostRoute(db *sqlx.DB, router *mux.Router, telegram service.TelegramService) {
	repo := repository.NewPostRepository(db)
	s := service.NewPostService(repo, 5*time.Second, telegram)
	c := controllers.NewPostController(s)

	router.HandleFunc("/create", c.Create).Methods("POST")
	router.HandleFunc("/{id}", c.Delete).Methods("DELETE")
	router.HandleFunc("/{id}", c.Find).Methods("GET")
	router.HandleFunc("/", c.GetAll).Methods("GET")
}
