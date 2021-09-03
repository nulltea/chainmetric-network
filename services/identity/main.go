package main

import (
	"github.com/timoth-y/chainmetric-network/services/identity/api/middleware"
	"github.com/timoth-y/chainmetric-network/services/identity/api/rpc"
	"github.com/timoth-y/chainmetric-network/services/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/services/identity/usecase/privileges"
	"github.com/timoth-y/chainmetric-network/shared/core"
	"github.com/timoth-y/chainmetric-network/shared/server"
	"github.com/timoth-y/chainmetric-network/shared/utils"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	privileges.Init()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
	utils.MustExecute(func() error {
		return server.Init(
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
				rpc.RegisterAdminService,
				rpc.RegisterUserService,
			),
		)
	}, "failed to initialize server")
}

func main() {
	if err := server.Serve(":8080"); err != nil {
		core.Logger.Fatal(err)
	}
}
