package main

import (
	"log"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	pkgDb "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/metrics"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/proto"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/delivery/grpc/server"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/search/repository/postgres"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8089"

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("can't listet port", zap.Error(err))
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
		logger.Error("Failed to setup metrics", zap.Error(err))
		os.Exit(1)
	}

	mw := middleware.NewGRPCMetricsMiddleware(ms)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.MetricsInterceptor),
	)

	proto.RegisterSearchEngineServer(server, searchServ)

	go metrics.ServePrometheusHTTP("0.0.0.0:9004")

	logger.Info("Starting search service")

	err = server.Serve(lis)
	if err != nil {
		logger.Error("Failed to start search server, ", zap.Error(err))
		return
	}
}
