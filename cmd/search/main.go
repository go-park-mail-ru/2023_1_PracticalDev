package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/consul"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/repository/postgres"
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
	config.DefaultGRPCSearchConfig()
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/src/configs/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read configuration", zap.Error(err))
	}

	// creating shortener service
	db, err := pkgDb.New(logger)
	if err != nil {
		os.Exit(1)
	}

	searchRepo := rep.NewRepository(db, logger)
	searchServ := serv.NewSearchServer(searchRepo)

	// setting up metrics
	ms := metrics.NewPrometheusMetrics("search")
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

	proto.RegisterSearchEngineServer(server, searchServ)
	lis, err := net.Listen("tcp", config.GetGRPCAddr())
	if err != nil {
		logger.Error("can't listet port", zap.Error(err))
		return
	}

	// starting server
	go metrics.ServePrometheusHTTP(viper.GetString(config.MetricsConfig.Addr))

	logger.Info("Starting search service")
	go func() {
		err = server.Serve(lis)
		if err != nil {
			logger.Error("Failed to start search server, ", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	server.GracefulStop()
}
