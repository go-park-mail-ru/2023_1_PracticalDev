package main

import (
	"context"
	"net/http"
	"os"

	pinsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/delivery/http"
	pinsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/repository/postgres"
	pinsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/service"

	_boardsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/delivery/http"
	_boardsRepo "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/repository/postgres"
	_boardsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/service"

	imagesRepo "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/repository/s3"
	_imagesServ "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/service"

	profileDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/delivery/http"
	_profileRepo "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/repository/postgres"
	_profileServ "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/service"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	_db "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"

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

	bucket, err := imagesRepo.NewS3Repository(logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := _imagesServ.NewS3Service(bucket)

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

	pinsRepo := pinsRepository.NewRepository(db, imagesServ, logger)
	pinsServ := pinsService.NewService(pinsRepo)

	profileRepo := _profileRepo.NewPostgresRepository(db, logger)
	profileServ := _profileServ.NewProfileService(profileRepo)

	auth.RegisterHandlers(mux, logger, authServ)
	users.RegisterHandlers(mux, logger, authorizer, users.NewService(users.NewRepository(db, logger)))
	profileDelivery.RegisterHandlers(mux, logger, authorizer, profileServ)
	_boardsDelivery.RegisterHandlers(mux, logger, authorizer, boardsAccessChecker, boardsServ)
	ping.RegisterHandlers(mux, logger)
	pinsDelivery.RegisterHandlers(mux, logger, authorizer, middleware.NewAccessChecker(pinsServ), pinsServ)

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
