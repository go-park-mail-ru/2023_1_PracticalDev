package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/repository/postgres"
	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/consul"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/redis"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

	// reading configs
	config.DefaultGRPCAuthConfig()
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/src/configs/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read configuration", zap.Error(err))
	}

	// creating auth service
	db, err := pkgDb.New(logger)
	if err != nil {
		logger.Error("Failed to create new postgres client", zap.Error(err))
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

	// setting up metrics
	ms := metrics.NewPrometheusMetrics("auth")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error("Failed to setup metrics", zap.Error(err))
		os.Exit(1)
	}
	mw := middleware.NewGRPCMetricsMiddleware(ms)

	// registration in consul
	cnsl, err := consul.NewConsulClient()
	if err != nil {
		logger.Error("Failed to connect to consul", zap.Error(err))
		os.Exit(1)
	}
	err = consul.RegisterService(cnsl)
	if err != nil {
		logger.Error("Failed to add service to consul", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("registered in consul")
	defer func() {
		err := cnsl.Agent().ServiceDeregister(viper.GetString(config.GrpcConfig.ServiceName) + "_" + config.GetConsulAddr())
		if err != nil {
			logger.Error("Failed to remove service from consul", zap.Error(err))
		}
		logger.Info("deregistered in consul")
	}()

	// creating grpc server
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
		grpc.MaxRecvMsgSize(viper.GetInt(config.GrpcConfig.MessageSize)*1024*1024),
	)
	proto.RegisterAuthenficatorServer(server, authServ)
	lis, err := net.Listen("tcp", config.GetGRPCAddr())
	if err != nil {
		logger.Error("can't listet addr", zap.Error(err))
		return
	}

	// starting server
	go metrics.ServePrometheusHTTP(viper.GetString(config.MetricsConfig.Addr))

	logger.Info("Starting auth service")
	go func() {
		err = server.Serve(lis)
		if err != nil {
			logger.Error("Failed to start auth server, ", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	server.GracefulStop()
}
