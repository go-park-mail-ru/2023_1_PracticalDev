package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/server"
	delHTTP "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/http"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/repository/mongo"
	service "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/service"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8090"

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("can't listet port", zap.Error(err))
		return
	}

	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI("mongodb://mongo:27017").
		SetAuth(options.Credential{
			Username: os.Getenv("MONGO_ROOT_USER"),
			Password: os.Getenv("MONGO_ROOT_PASSWORD"),
		}))
	if err != nil {
		logger.Error("failed to connect to mongo db", zap.Error(err))
		os.Exit(1)
	}
	defer client.Disconnect(context.Background()) //nolint
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Error("failed to ping mongo db", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("connected to mongo db successfully")

	shortenerRep := rep.NewShortenerRepository(client.Database("url_shortener"), logger)
	shortenerServ := serv.NewShortenerServer(shortenerRep)

	ms := metrics.NewPrometheusMetrics("shortener")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error("Failed to setup metrics", zap.Error(err))
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterShortenerServer(server, shortenerServ)

	logger.Info("Starting shortener grpc service")
	go func() {
		err = server.Serve(lis)
		if err != nil {
			logger.Error("Failed to start shortener grpc server, ", zap.Error(err))
			return
		}
	}()

	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	shortenerService := service.NewShortenerService(shortenerRep)

	metricsMiddleware := middleware.NewHttpMetricsMiddleware(ms)

	delHTTP.RegisterGetHandler(mux, logger, shortenerService, metricsMiddleware)

	httpServer := http.Server{
		Addr:    "0.0.0.0:8091",
		Handler: mux,
	}

	logger.Info("Starting http metrics...")

	go metrics.ServePrometheusHTTP("0.0.0.0:9001")

	logger.Info("Starting http server...")
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
