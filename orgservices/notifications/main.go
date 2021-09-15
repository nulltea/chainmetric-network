package main

import (
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/api/rpc"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/utils"
)

func init() {
	core.Init()
	utils.MustExecute(func() error {
		return server.Init(
			server.WithUnaryMiddlewares(
				middleware.JWTForUnaryGRPC(),
				middleware.AuthForUnaryGRPC(),
			),
			server.WithStreamMiddlewares(
				middleware.JWTForStreamGRPC(),
				middleware.AuthForStreamGRPC(),
			),
			server.WithServiceRegistrar(
				rpc.RegisterSubscriberService,
			),
		)
	}, "failed to initialize server")
}

func main() {
	utils.MustExecute(func() error {
		return server.Serve(":8080")
	}, "failed to initialize gRPC server")
}
