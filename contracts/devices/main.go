package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-network/shared/core"
	"github.com/ztrue/shutdown"
)

func init() {
	core.InitCore()
}

func main() {
	go core.BootstrapContract(NewDevicesContact())

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
