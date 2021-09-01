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

func WithUnaryMiddlewares(middlewares ...grpc.UnaryServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.UnaryInterceptors = append(args.UnaryInterceptors, middlewares...)
	}
}

func WithStreamMiddlewares(middlewares ...grpc.StreamServerInterceptor) Option {
	return func(args *gRPCArgsStub) {
		args.StreamInterceptors = append(args.StreamInterceptors, middlewares...)
	}
}

func WithServiceRegistrar(registrar func(server *grpc.Server)) Option {
	return func(args *gRPCArgsStub) {
		args.ServicesRegistrars = append(args.ServicesRegistrars, registrar)
	}
}
