package main

import (
	"github.com/timoth-y/chainmetric-contracts/shared"
)

func init() {
	shared.InitCore()
}

func main() {
	shared.BootstrapContract(NewReadingsContract())
}
