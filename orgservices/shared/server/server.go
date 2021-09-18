package server

import (
	"net"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var server *grpc.Server

// Init performs initialization of the gRPC server.
func Init(options ...Option) error {
	var (
		args = &gRPCArgsStub{
			unaryInterceptors: []grpc.UnaryServerInterceptor{
				tags.UnaryServerInterceptor(),
				recovery.UnaryServerInterceptor(),
			},

			streamInterceptors: []grpc.StreamServerInterceptor{
				tags.StreamServerInterceptor(),
				recovery.StreamServerInterceptor(),
			},
		}

		certPath = viper.GetString("api.grpc_tls_cert")
		keyPath  = viper.GetString("api.grpc_tls_key")
	)

	for i := range options {
		options[i](args)
	}

	if entry := logrus.NewEntry(args.logger); entry.Logger != nil {
		grpclogrus.ReplaceGrpcLogger(entry)

		args.unaryInterceptors = append(args.unaryInterceptors, grpclogrus.UnaryServerInterceptor(entry))
		args.streamInterceptors = append(args.streamInterceptors, grpclogrus.StreamServerInterceptor(entry))
	}

	tls, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return err
	}

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(args.unaryInterceptors...),
		grpc.ChainStreamInterceptor(args.streamInterceptors...),
		grpc.Creds(tls),
	)

	for i := range args.servicesRegistrars {
		args.servicesRegistrars[i](server)
	}

	return nil
}

// Serve starts gRPC server on given `addr`.
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
