package main

import (
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/server"
	"github.com/timoth-y/chainmetric-contracts/shared/utils"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/rpc"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
	utils.MustExecute(func() error {
		return server.Init(
			rpc.WithIdentityService,
		)
	}, "failed to initialize server")
}

func main() {
	if err := server.Serve(":8080"); err != nil {
		core.Logger.Fatal(err)
	}
}
