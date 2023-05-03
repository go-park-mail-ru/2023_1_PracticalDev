package main

import (
	"log"
	"net"
	"os"
	"strconv"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/repository/s3"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/service/grpc"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zaplogger.New()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	_, err = strconv.Atoi(port)
	if port == "" || err != nil {
		port = "8088"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("can't listet port", err)
		os.Exit(1)
	}

	bucket, err := rep.NewS3Repository(logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := serv.NewS3Service(bucket)

	ms := metrics.NewPrometheusMetrics("images")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterImageUploaderServer(server, imagesServ)

	go metrics.ServePrometheusHTTP("0.0.0.0:9002")

	logger.Info("Starting images service")

	err = server.Serve(lis)
	if err != nil {
		logger.Error("Failed to start image server, ", err.Error())
		return
	}
}
