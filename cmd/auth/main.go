package main

import (
	"context"
	"log"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/repository/postgres"
	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8087"

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
	logger := zap.New(consoleCore)
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("can't listet port", zap.Error(err))
		return
	}

	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	rdb, err := redis.NewRedisClient(logger, ctx)
	if err != nil {
		logger.Error("Failed to create new redis client", zap.Error(err))
		os.Exit(1)
	}

	authRepo := rep.NewRepository(db, rdb, ctx, logger)
	authServ := serv.NewAuthServer(authRepo)

	ms := metrics.NewPrometheusMetrics("auth")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error("Failed to setup metrics", zap.Error(err))
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterAuthenficatorServer(server, authServ)

	go metrics.ServePrometheusHTTP("0.0.0.0:9003")

	logger.Info("Starting auth service")

	err = server.Serve(lis)
	if err != nil {
		logger.Error("Failed to start auth server, ", zap.Error(err))
		return
	}
}
