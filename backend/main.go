package main

import (
	"blog-app/internal/api/route"
	"blog-app/internal/setup"
	"blog-app/pkg/service"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

type App struct {
	Env      *setup.Env
	DB       *sqlx.DB
	Telegram *service.TelegramService
}

func main() {
	r := mux.NewRouter()
	app := App{}
	app.Env = setup.NewEnv()
	app.DB = setup.ConnectToDB(app.Env)
	bot, err := gotgbot.NewBot(app.Env.TGToken, nil)
	if err != nil {
		log.Fatal(err)
	}
	app.Telegram = service.NewTelegramService(bot, app.Env.ChannelID)

	postsGroup := r.PathPrefix("/posts").Subrouter()
	route.NewPostRoute(app.DB, postsGroup, app.Telegram)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
