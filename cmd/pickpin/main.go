package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/client"
	authDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/http"

	likesDelivery "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/delivery/http"
	likesRepository "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/repository/postgres"
	likesService "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes/service"

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

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/ping"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
)

func main() {
	// Zap logger configuration
	consoleCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Zap logger
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	logger := zap.New(consoleCore, zap.AddCaller())
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	imagesConn, err := grpc.Dial(
		"images:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("cant connect to images service")
		os.Exit(1)
	}
	imagesServ := imagesService.NewImageUploaderClient(imagesConn)

	authConn, err := grpc.Dial(
		"auth:8087",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("cant connect to images service")
		os.Exit(1)
	}
	authServ := authService.NewAuthenficatorClient(authConn)

	searchConn, err := grpc.Dial(
		"search:8089",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("cant connect to search service")
		os.Exit(1)
	}

	mt := metrics.NewPrometheusMetrics("pickpin")
	err = mt.SetupMetrics()
	metricsMiddleware := middleware.NewHttpMetricsMiddleware(mt)

	if err != nil {
		logger.Error("failed to setup prometheus", zap.Error(err))
		os.Exit(1)
	}

	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	likesRepo := likesRepository.NewRepository(db, logger)
	likesServ := likesService.NewService(likesRepo)

	token := tokens.NewHMACHashToken(config.Get("CSRF_TOKEN_SECRET"))
	CSRFMiddleware := middleware.NewCSRFMiddleware(token, logger)
	authorizer := middleware.NewAuthorizer(authServ, logger)

	pinsRepo := pinsRepository.NewRepository(db, imagesServ, logger)
	pinsServ := pinsService.NewService(pinsRepo)

	searchServ := searchService.NewSearchClient(searchConn, pinsServ)

	boardsRepo := boardsRepository.NewPostgresRepository(db, logger)
	boardsServ := boardsService.NewBoardsService(boardsRepo, pinsServ)
	boardsAccessChecker := middleware.NewAccessChecker(boardsServ)

	usersRepo := usersRepository.NewRepository(db, logger)
	usersServ := usersService.NewService(usersRepo)

	profileRepo := profileRepository.NewPostgresRepository(db, imagesServ, logger)
	profileServ := profileService.NewProfileService(profileRepo)

	followingsRepo := followingsRepository.NewRepository(db, logger)
	followingsServ := followingsService.NewService(followingsRepo)

	chatsRepo := chatsRepository.NewRepository(db, logger)
	chatsServ := chatsService.NewService(chatsRepo)

	authDelivery.RegisterHandlers(mux, logger, authServ, token, metricsMiddleware)
	likesDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, likesServ, metricsMiddleware)
	usersDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, usersServ, metricsMiddleware)
	profileDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, profileServ, metricsMiddleware)
	followingsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, followingsServ, metricsMiddleware)
	boardsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, boardsAccessChecker, boardsServ, metricsMiddleware)
	pinsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, middleware.NewAccessChecker(pinsServ), pinsServ, metricsMiddleware)
	chatsDelivery.RegisterHandlers(mux, logger, authorizer, CSRFMiddleware, chatsServ)
	ping.RegisterHandlers(mux, logger)
	searchDelivery.RegisterHandlers(mux, logger, authorizer, searchServ)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	logger.Info("Starting metrics...")

	go metrics.ServePrometheusHTTP("0.0.0.0:9001")

	logger.Info("Starting server...")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
