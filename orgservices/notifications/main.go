package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/api/rpc"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/usecase/eventproxy"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/ztrue/shutdown"
)

func init() {
	core.Init()
	shutdown.Add(eventproxy.Stop)
}

func main() {
	eventproxy.Start()

	go core.BootstrapGRPCServer(8080,
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

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
