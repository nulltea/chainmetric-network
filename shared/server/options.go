package server

import "google.golang.org/grpc"

type (
	// Option defines parameter for configuring grpc.Server.
	Option func(server *grpc.Server)
)
