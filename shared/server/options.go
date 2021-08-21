package server

import "google.golang.org/grpc"

type (
	gRPCArgsStub struct {
		UnaryInterceptors  []grpc.UnaryServerInterceptor
		StreamInterceptors []grpc.StreamServerInterceptor
		ServicesRegistrars []func(server *grpc.Server)
	}

	// Option defines parameter for configuring grpc.Server.
	Option func(args *gRPCArgsStub)
)

func WithUnaryMiddleware(middleware grpc.UnaryServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.UnaryInterceptors = append(args.UnaryInterceptors, middleware)
	}
}

func WithStreamMiddleware(middleware grpc.StreamServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.StreamInterceptors = append(args.StreamInterceptors, middleware)
	}
}

func WithServiceRegistrar(registrar func(server *grpc.Server)) Option {
	return func(args *gRPCArgsStub) {
		args.ServicesRegistrars = append(args.ServicesRegistrars, registrar)
	}
}
