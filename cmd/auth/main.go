package main

import (
	"context"
	"log"
	"net"
	"os"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/repository/postgres"
	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8087"

func main() {
	logger, err := zaplogger.New()

	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("can't listet port", err)
		return
	}

	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	rdb, err := redis.NewRedisClient(logger, ctx)
	if err != nil {
		logger.Warn(err)
		os.Exit(1)
	}

	authRepo := rep.NewRepository(db, rdb, ctx, logger)
	authServ := serv.NewAuthServer(authRepo)

	ms := metrics.NewPrometheusMetrics("auth")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterAuthenficatorServer(server, authServ)

	go metrics.ServePrometheusHTTP("0.0.0.0:9003")

	logger.Info("Starting auth service")

	server.Serve(lis)
}
