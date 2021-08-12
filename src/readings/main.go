package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/ztrue/shutdown"
)

func init() {
	core.InitCore()
	core.InitLevelDB()
}

func main() {
	go core.BootstrapContract(NewReadingsContract())

	shutdown.Add(core.CloseLevelDB)
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
