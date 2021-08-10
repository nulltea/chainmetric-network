package main

import (
	"syscall"

	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/ztrue/shutdown"
)

func init() {
	shared.InitCore()
	shared.InitLevelDB()
}

func main() {
	go shared.BootstrapContract(NewReadingsContract())

	shutdown.Add(shared.CloseLevelDB)
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
