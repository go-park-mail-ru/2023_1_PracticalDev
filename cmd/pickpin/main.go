package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/posts"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.New()

	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	rdb, err := redis.NewRedisClient(logger, ctx)
	if err != nil {
		logger.Warn(err)
		os.Exit(1)
	}

	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	authServ := auth.NewService(auth.NewRepository(db, rdb, ctx, logger))
	authorizer := middleware.NewAuthorizer(authServ)

	users.RegisterHandlers(mux, logger, authorizer, users.NewService(users.NewRepository(db, logger)))
	auth.RegisterHandlers(mux, logger, authServ)
	posts.RegisterHandlers(mux, logger, authorizer, posts.NewService(posts.NewRepository(db, logger)))
	ping.RegisterHandlers(mux, logger)

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
