package main

import "github.com/timoth-y/chainmetric-contracts/shared"

func init() {
	shared.InitLogger()
}

func main() {
	shared.BootstrapContract(NewDevicesContact())
}
