package server

import (
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	logger *logrus.Logger
	server *grpc.Server
)

func Init(options ...Option) {
	var (
		unaryInterceptors = []grpc.UnaryServerInterceptor{
			tags.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(),
			grpclogrus.UnaryServerInterceptor(initLogger()),
		}

		streamInterceptors = []grpc.StreamServerInterceptor{
			tags.StreamServerInterceptor(),
			recovery.StreamServerInterceptor(),
			grpclogrus.StreamServerInterceptor(initLogger()),
		}
	)

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	for i := range options {
		options[i](server)
	}

	reflection.Register(server)
}

func Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "failed to create listener")
	}

	if err = server.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func initLogger() *logrus.Entry {
	logger = logrus.New()
	loggerEntry := logrus.NewEntry(logger)
	grpclogrus.ReplaceGrpcLogger(loggerEntry)

	return loggerEntry
}
