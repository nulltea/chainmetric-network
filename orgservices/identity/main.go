package main

import (
	"github.com/timoth-y/chainmetric-network/orgservices/identity/api/rpc"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/privileges"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/utils"
)

func init() {
	core.Init()
	privileges.Init()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
}

func main() {
	core.BootstrapGRPCServer(8080,
		server.WithUnaryMiddlewares(
			middleware.JWTForUnaryGRPC(
				"UserService/register",
				"AccessService/requestFabricCredentials",
				"AccessService/authWithSigningIdentity",
			),
			middleware.AuthForUnaryGRPC("UserService/pingAccountStatus"),
		),
		server.WithStreamMiddlewares(
			middleware.JWTForStreamGRPC(),
			middleware.AuthForStreamGRPC(),
		),
		server.WithServiceRegistrar(
			rpc.RegisterAccessService,
			rpc.RegisterUserService,
			rpc.RegisterAdminService,
		),
	)
}
