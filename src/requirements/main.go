package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/ztrue/shutdown"
)

func init() {
	core.InitCore()
}

func main() {
	go core.BootstrapContract(NewRequirementsContract())

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
