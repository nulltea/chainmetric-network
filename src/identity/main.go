package main

import (
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/server"
	"github.com/timoth-y/chainmetric-contracts/shared/utils"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/middleware"
	"github.com/timoth-y/chainmetric-contracts/src/identity/api/rpc"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/identity"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
	utils.MustExecute(func() error {
		return server.Init(
			server.WithUnaryMiddleware(middleware.JWTAuthUnaryInterceptor(
				"IdentityService/register", "AuthService/authenticate",
			)),
			server.WithStreamMiddleware(middleware.JWTAuthStreamInterceptor()),
			server.WithServiceRegistrar(rpc.RegisterIdentityService),
			server.WithServiceRegistrar(rpc.RegisterAuthService),
		)
	}, "failed to initialize server")
}

func main() {
	if err := server.Serve(":8080"); err != nil {
		core.Logger.Fatal(err)
	}
}
