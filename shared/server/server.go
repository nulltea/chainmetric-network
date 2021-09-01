package server

import (
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var server *grpc.Server

func Init(options ...Option) error {
	var (
		args = &gRPCArgsStub{
			UnaryInterceptors: []grpc.UnaryServerInterceptor{
				tags.UnaryServerInterceptor(),
				recovery.UnaryServerInterceptor(),
				grpclogrus.UnaryServerInterceptor(initLogger()),
			},

			StreamInterceptors: []grpc.StreamServerInterceptor{
				tags.StreamServerInterceptor(),
				recovery.StreamServerInterceptor(),
				grpclogrus.StreamServerInterceptor(initLogger()),
			},
		}

		certPath = viper.GetString("api.grpc_tls_cert")
		keyPath  = viper.GetString("api.grpc_tls_key")
	)

	for i := range options {
		options[i](args)
	}

	tls, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return err
	}

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(args.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(args.StreamInterceptors...),
		grpc.Creds(tls),
	)

	for i := range args.ServicesRegistrars {
		args.ServicesRegistrars[i](server)
	}

	return nil
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
	loggerEntry := logrus.NewEntry(core.Logrus)
	grpclogrus.ReplaceGrpcLogger(loggerEntry)

	return loggerEntry
}
