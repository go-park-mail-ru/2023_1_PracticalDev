package main

import (
	"log"
	"net"
	"os"
	"strconv"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/repository/s3"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/service/grpc"
	zaplogger "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log/zap"
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

	server := grpc.NewServer()
	bucket, err := rep.NewS3Repository(logger)
	if err != nil {
		os.Exit(1)
	}
	imagesServ := serv.NewS3Service(bucket)

	proto.RegisterImageUploaderServer(server, imagesServ)

	server.Serve(lis)
}
