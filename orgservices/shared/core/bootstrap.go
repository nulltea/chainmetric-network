package core

import (
	"fmt"

	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/utils"
)

// BootstrapGRPCServer performs start sequence of the gRPC server.
func BootstrapGRPCServer(port int, options ...server.Option) {
	utils.MustExecute(func() error {
		return server.Init(options...)
	}, "failed to initialize server")

	utils.MustExecute(func() error {
		return server.Serve(fmt.Sprintf(":%d", port))
	}, "failed to initialize gRPC server")
}
