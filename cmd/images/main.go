package main

import (
	"net"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"

	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/proto"
	rep "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/repository/s3"
	serv "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images/server/service/grpc"
	"google.golang.org/grpc"
)

var port = "0.0.0.0:8088"

func main() {
	logger := log.New()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("can't listet port", err)
		return
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
