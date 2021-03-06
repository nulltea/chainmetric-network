package main

import "github.com/timoth-y/iot-blockchain-contracts/shared"

func init() {
	shared.InitLogger()
}

func main() {
	shared.BootstrapContract(NewAssetsContact())
}
