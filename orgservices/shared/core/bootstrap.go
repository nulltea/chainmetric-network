package core

import (
	"fmt"

	"github.com/timoth-y/chainmetric-network/orgservices/shared/server"
)

// BootstrapGRPCServer performs start sequence of the gRPC server.
func BootstrapGRPCServer(port int, options ...server.Option) {
	if err := server.Init(options...); err != nil {
		Logrus.WithError(err).Fatalln("failed to initialize server")
	}

	if err := server.Serve(fmt.Sprintf(":%d", port)); err != nil {
		Logrus.WithError(err).Fatalln("failed to start gRPC server")
	}
}
