package main

import (
	"context"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	authDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/http"
	authRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/repository/postgres"
	authService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/service"

	likesDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/delivery/http"
	likesRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/repository/postgres"
	likesService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/service"

	pinsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/delivery/http"
	pinsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/repository/postgres"
	pinsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/service"

	boardsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/delivery/http"
	boardsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/repository/postgres"
	boardsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/service"

	imagesRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/repository/s3"
	imagesService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/service"

	usersDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/delivery/http"
	usersRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/repository/postgres"
	usersService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/service"

	profileDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/delivery/http"
	profileRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/repository/postgres"
	profileService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/service"

	followingsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/delivery/http"
	followingsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/repository/postgres"
	followingsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/service"

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
)

func main() {
	logger := log.New()

	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	bucket, err := imagesRepository.NewS3Repository(logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := imagesService.NewS3Service(bucket)

	ctx := context.Background()
	rdb, err := redis.NewRedisClient(logger, ctx)
	if err != nil {
		logger.Warn(err)
		os.Exit(1)
	}

	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	boardsRepo := boardsRepository.NewPostgresRepository(db, logger)
	boardsServ := boardsService.NewBoardsService(boardsRepo)
	boardsAccessChecker := middleware.NewAccessChecker(boardsServ)

	likesRepo := likesRepository.NewRepository(db, logger)
	likesServ := likesService.NewService(likesRepo)

	token := tokens.NewHMACHashToken(config.Get("CSRF_TOKEN_SECRET"))
	authRepo := authRepository.NewRepository(db, rdb, ctx, logger)
	authServ := authService.NewService(authRepo)
	authorizer := middleware.NewAuthorizer(authServ, token, logger)

	pinsRepo := pinsRepository.NewRepository(db, imagesServ, logger)
	pinsServ := pinsService.NewService(pinsRepo)

	usersRepo := usersRepository.NewRepository(db, logger)
	usersServ := usersService.NewService(usersRepo)

	profileRepo := profileRepository.NewPostgresRepository(db, imagesServ, logger)
	profileServ := profileService.NewProfileService(profileRepo)

	followingsRepo := followingsRepository.NewRepository(db, logger)
	followingsServ := followingsService.NewService(followingsRepo)

	authDelivery.RegisterHandlers(mux, logger, authServ, token)
	likesDelivery.RegisterHandlers(mux, logger, authorizer, likesServ)
	usersDelivery.RegisterHandlers(mux, logger, authorizer, usersServ)
	profileDelivery.RegisterHandlers(mux, logger, authorizer, profileServ)
	followingsDelivery.RegisterHandlers(mux, logger, authorizer, followingsServ)
	boardsDelivery.RegisterHandlers(mux, logger, authorizer, boardsAccessChecker, boardsServ)
	pinsDelivery.RegisterHandlers(mux, logger, authorizer, middleware.NewAccessChecker(pinsServ), pinsServ)
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
