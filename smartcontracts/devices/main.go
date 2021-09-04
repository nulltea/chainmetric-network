package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/ztrue/shutdown"
)

func init() {
	core.InitCore()
}

func main() {
	go core.BootstrapChaincodeServer(NewDevicesContact())

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
