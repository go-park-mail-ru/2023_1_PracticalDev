package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	authService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/client"
	authDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/http"

	likesDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/delivery/http"
	likesRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/repository/postgres"
	likesService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/service"

	notificationsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/delivery/http"
	notificationsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/repository/postgres"
	notificationsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/service"

	pinsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/delivery/http"
	pinsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/repository/postgres"
	pinsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/service"

	boardsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/delivery/http"
	boardsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/repository/postgres"
	boardsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards/service"

	imagesService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/client"

	usersDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/delivery/http"
	usersRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/repository/postgres"
	usersService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/service"

	profileDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/delivery/http"
	profileRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/repository/postgres"
	profileService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/service"

	followingsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/delivery/http"
	followingsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/repository/postgres"
	followingsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/service"

	searchService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/client"
	searchDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/http"

	chatsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats/delivery/http"
	chatsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats/repository/postgres"
	chatsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats/service"

	commentsDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments/delivery/http"
	commentsRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments/repository/postgres"
	commentsService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments/service"

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	shortenerService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/client"
	shortenerDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/consul"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/resolvers"
)

func main() {
	// setting up logger
	logger := zaplogger.NewDefaultZapProdLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	// load config
	config.DefaultPickPinConfig()
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/src/configs/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read configuration", zap.Error(err))
	}

	// getting consul client
	cnsl, err := consul.NewConsulClient()
	if err != nil {
		logger.Error("failed to connect to consul", zap.Error(err))
	}

	ctx, cancelDiscovery := context.WithCancel(context.Background())

	imagesConn, err := resolvers.NewGRPCConnWithResolver(ctx, cnsl, "image", logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := imagesService.NewImageUploaderClient(imagesConn)

	authConn, err := resolvers.NewGRPCConnWithResolver(ctx, cnsl, "auth", logger)
	if err != nil {
		os.Exit(1)
	}
	authServ := authService.NewAuthenficatorClient(authConn)

	searchConn, err := resolvers.NewGRPCConnWithResolver(ctx, cnsl, "search", logger)
	if err != nil {
		os.Exit(1)
	}

	shortenerConn, err := resolvers.NewGRPCConnWithResolver(ctx, cnsl, "shortener", logger)
	if err != nil {
		os.Exit(1)
	}

	// setting up metrics
	mt := metrics.NewPrometheusMetrics("pickpin")
	err = mt.SetupMetrics()
	metricsMiddleware := middleware.NewHttpMetricsMiddleware(mt)

	if err != nil {
		logger.Error("failed to setup prometheus", zap.Error(err))
		os.Exit(1)
	}

	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	// creating services
	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	notificationsRepo := notificationsRepository.NewRepository(db, logger)
	notificationsServ := notificationsService.NewService(notificationsRepo, logger)

	followingsRepo := followingsRepository.NewRepository(db, logger)
	followingsServ := followingsService.NewService(followingsRepo, notificationsServ)

	pinsRepo := pinsRepository.NewRepository(db, imagesServ, logger)
	pinsServ := pinsService.NewService(pinsRepo, notificationsServ, followingsRepo)

	likesRepo := likesRepository.NewRepository(db, logger)
	likesServ := likesService.NewService(likesRepo, notificationsServ, pinsRepo, logger)

	token := tokens.NewHMACHashToken(viper.GetString(config.CSRFConfig.Token))
	CSRFMiddleware := middleware.NewCSRFMiddleware(token, logger)
	authorizer := middleware.NewAuthorizer(authServ, logger)

	searchServ := searchService.NewSearchClient(searchConn, pinsServ)

	boardsRepo := boardsRepository.NewPostgresRepository(db, logger)
	boardsServ := boardsService.NewBoardsService(boardsRepo, pinsServ)
	boardsAccessChecker := middleware.NewAccessChecker(boardsServ)

	usersRepo := usersRepository.NewRepository(db, logger)
	usersServ := usersService.NewService(usersRepo)

	profileRepo := profileRepository.NewPostgresRepository(db, imagesServ, logger)
	profileServ := profileService.NewProfileService(profileRepo)

	chatsRepo := chatsRepository.NewRepository(db, logger)
	chatsServ := chatsService.NewService(chatsRepo)

	commentsRepo := commentsRepository.NewRepository(db, logger)
	commentsServ := commentsService.NewService(commentsRepo, notificationsServ, pinsRepo)

	shortServ := shortenerService.NewShortenerClient(shortenerConn)

	authDelivery.RegisterHandlers(mux, logger, authServ, token, metricsMiddleware)
	likesDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, likesServ, metricsMiddleware)
	usersDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, usersServ, metricsMiddleware)
	profileDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, profileServ, metricsMiddleware)
	followingsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, followingsServ, metricsMiddleware)
	boardsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, boardsAccessChecker, boardsServ, metricsMiddleware)
	pinsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, middleware.NewAccessChecker(pinsServ), pinsServ, metricsMiddleware)
	chatsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, chatsServ)
	commentsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, commentsServ, metricsMiddleware)
	ping.RegisterHandlers(mux, logger)
	searchDelivery.RegisterHandlers(mux, logger, authorizer, searchServ)
	shortenerDelivery.RegisterPostHandler(mux, logger, authorizer, CSRFMiddleware, shortServ, metricsMiddleware)
	notificationsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, notificationsServ,
		metricsMiddleware)

	// setting up server
	server := http.Server{
		Addr:    viper.GetString(config.HttpConfig.Addr),
		Handler: mux,
	}

	logger.Info("Starting metrics...")

	go metrics.ServePrometheusHTTP(viper.GetString(config.MetricsConfig.Addr))

	logger.Info("Starting server...")
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancelDiscovery()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error("failed to gracefull shutdown http server", zap.Error(err))
	}
}
