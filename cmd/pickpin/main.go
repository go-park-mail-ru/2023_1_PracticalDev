package main

import (
	"context"

	_boardsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/delivery/http"
	_boardsRepo "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/repository/postgres"
	_boardsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/service"

	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	_db "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/posts"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.New()

	db, err := _db.New(logger)
	if err != nil {
		os.Exit(1)
	}

	bucket, err := images.NewRepository(logger)
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

	boardsRepo := _boardsRepo.NewPostgresRepository(db, logger)
	boardsServ := _boardsService.NewBoardsService(boardsRepo)
	boardsAccessChecker := middleware.NewAccessChecker(boardsServ)

	authServ := auth.NewService(auth.NewRepository(db, rdb, ctx, logger))
	authorizer := middleware.NewAuthorizer(authServ)

	pinsS3 := images.NewS3Service(bucket)
	pinsRepo := pins.NewRepository(db, pinsS3, logger)
	pinsServ := pins.NewService(pinsRepo)

	auth.RegisterHandlers(mux, logger, authServ)
	users.RegisterHandlers(mux, logger, authorizer, users.NewService(users.NewRepository(db, logger)))
	posts.RegisterHandlers(mux, logger, authorizer, posts.NewService(posts.NewRepository(db, logger)))
	_boardsDelivery.RegisterHandlers(mux, logger, authorizer, boardsAccessChecker, boardsServ)
	ping.RegisterHandlers(mux, logger)
	pins.RegisterHandlers(mux, logger, authorizer, middleware.NewAccessChecker(pinsServ), pinsServ)

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
