package main

import (
	"log"
	"net"
	"os"

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/repository/postgres"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8089"

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

	searchRepo := rep.NewRepository(db, logger)
	searchServ := serv.NewSearchServer(searchRepo)

	ms := metrics.NewPrometheusMetrics("search")
	err = ms.SetupMetrics()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterSearchEngineServer(server, searchServ)

	go metrics.ServePrometheusHTTP("0.0.0.0:9004")

	logger.Info("Starting search service")

	server.Serve(lis)
}
