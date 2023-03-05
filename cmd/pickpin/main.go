package main

import (
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/posts"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.New()

	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}

	mux := httprouter.New()

	users.RegisterHandlers(mux, logger, users.NewService(users.NewRepository(db, logger)))
	posts.RegisterHandlers(mux, logger, posts.NewService(posts.NewRepository(db, logger)))

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	logger.Info("Starting server...")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server, ", err.Error())
	}
}
