package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/repository/s3"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/service/grpc"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/consul"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	config.DefaultGRPCImageConfig()
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/src/configs/")
	viper.ReadInConfig()

	// creating image service
	bucket, err := rep.NewS3Repository(logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := serv.NewS3Service(bucket)

	// setting up metrics
	ms := metrics.NewPrometheusMetrics("images")
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
	proto.RegisterImageUploaderServer(server, imagesServ)
	lis, err := net.Listen("tcp", config.GetGRPCAddr())
	if err != nil {
		logger.Error("can't listet addr", zap.Error(err))
		os.Exit(1)
	}

	// starting server
	go metrics.ServePrometheusHTTP(viper.GetString(config.MetricsConfig.Addr))

	logger.Info("Starting images service")
	go func() {
		err = server.Serve(lis)
		if err != nil {
			logger.Error("Failed to start image server, ", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	server.GracefulStop()
}
