package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/api/rpc"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/usecase/eventproxy"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/usecase/privileges"
	"github.com/ztrue/shutdown"
)

func init() {
	core.Init()
	privileges.Init()
	eventproxy.Init()
	shutdown.Add(eventproxy.Stop)
}

func main() {
	eventproxy.Start()

	go core.BootstrapGRPCServer(8080,
		server.WithUnaryMiddlewares(
			middleware.JWTForUnaryGRPC(),
			middleware.AuthForUnaryGRPC(),
			middleware.FirebaseForUnaryGRPC(true),
		),
		server.WithStreamMiddlewares(
			middleware.JWTForStreamGRPC(),
			middleware.AuthForStreamGRPC(),
			middleware.FirebaseForStreamGRPC(true),
		),
		server.WithServiceRegistrar(
			rpc.RegisterSubscriberService,
		),
		server.WithLogger(core.Logrus),
	)

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
