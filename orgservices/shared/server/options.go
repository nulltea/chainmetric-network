package server

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	gRPCArgsStub struct {
		unaryInterceptors    []grpc.UnaryServerInterceptor
		streamInterceptors   []grpc.StreamServerInterceptor
		servicesRegistrars   []RegistrarFunc
		transportCredentials credentials.TransportCredentials
		logger *logrus.Logger
	}

	// Option defines parameter for configuring grpc.Server.
	Option func(args *gRPCArgsStub)

	// RegistrarFunc is a function that registers proto service implementation against gRPC server.
	RegistrarFunc func(server *grpc.Server)
)

// WithUnaryMiddlewares can be used to pass grpc.UnaryServerInterceptor for server.
func WithUnaryMiddlewares(middlewares ...grpc.UnaryServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.unaryInterceptors = append(args.unaryInterceptors, middlewares...)
	}
}

// WithStreamMiddlewares can be used to pass grpc.StreamServerInterceptor for server.
func WithStreamMiddlewares(middlewares ...grpc.StreamServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.streamInterceptors = append(args.streamInterceptors, middlewares...)
	}
}

// WithServiceRegistrar can be used to register proto gRPC services against server.
func WithServiceRegistrar(registrars ...RegistrarFunc) Option {
	return func(args *gRPCArgsStub) {
		args.servicesRegistrars = append(args.servicesRegistrars, registrars...)
	}
}

// WithLogger can be used to pass logger for logging gRPC transactions.
func WithLogger(logger *logrus.Logger) Option {
	return func(args *gRPCArgsStub) {
		args.logger = logger
	}
}
