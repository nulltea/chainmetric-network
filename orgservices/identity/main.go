package main

import (
	"github.com/timoth-y/chainmetric-network/orgservices/identity/api/middleware"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/privileges"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	server2 "github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/utils"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	privileges.Init()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
	utils.MustExecute(func() error {
		return server2.Init(
			server2.WithUnaryMiddlewares(
				middleware.JWTForUnaryGRPC(
					"UserService/register",
					"AccessService/requestFabricCredentials",
					"AccessService/authWithSigningIdentity",
				),
				middleware.AuthForUnaryGRPC("UserService/pingAccountStatus"),
			),
			server2.WithStreamMiddlewares(
				middleware.JWTForStreamGRPC(),
				middleware.AuthForStreamGRPC(),
			),
			server2.WithServiceRegistrar(
			),
		)
	}, "failed to initialize server")
}

func main() {
	utils.MustExecute(func() error {
		return server2.Serve(":8080")
	}, "failed to initialize gRPC server")
}
