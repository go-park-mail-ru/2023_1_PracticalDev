package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/consul"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/mongo"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/server"
	delHTTP "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/http"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/repository/mongo"
	service "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/service"
	"github.com/julienschmidt/httprouter"
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
	config.DefaultGRPCShortenerConfig()
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/src/configs/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read configuration", zap.Error(err))
	}

	// creating shortener service
	client, err := mongo.NewMongoClient(logger)
	if err != nil {
		os.Exit(1)
	}
	defer client.Disconnect(context.Background()) //nolint

	shortenerRep := rep.NewShortenerRepository(client.Database("url_shortener"), logger)
	shortenerServ := serv.NewShortenerServer(shortenerRep)

	// setting up metrics
	ms := metrics.NewPrometheusMetrics("shortener")
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
	)
	proto.RegisterShortenerServer(server, shortenerServ)
	lis, err := net.Listen("tcp", config.GetGRPCAddr())
	if err != nil {
		logger.Error("can't listet port", zap.Error(err))
		return
	}
	logger.Info("Starting shortener grpc service")
	go func() {
		err = server.Serve(lis)
		if err != nil {
			logger.Error("Failed to start shortener grpc server, ", zap.Error(err))
			return
		}
	}()

	// setting up http server
	mux := httprouter.New()
	mux.GlobalOPTIONS = middleware.HandlerFuncLogger(middleware.OptionsHandler, logger)

	shortenerService := service.NewShortenerService(shortenerRep)

	metricsMiddleware := middleware.NewHttpMetricsMiddleware(ms)
	delHTTP.RegisterGetHandler(mux, logger, shortenerService, metricsMiddleware)

	httpServer := http.Server{
		Addr:    viper.GetString(config.HttpConfig.Addr),
		Handler: mux,
	}

	logger.Info("Starting metrics...")
	go metrics.ServePrometheusHTTP(viper.GetString(config.MetricsConfig.Addr))

	logger.Info("Starting http server...")
	go func() {
		err = httpServer.ListenAndServe()
		if err != nil {
			logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	server.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("failed to gracefull shutdown http server", zap.Error(err))
	}
}
