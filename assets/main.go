package main

//go:generate protoc --proto_path=model/ --go_out=paths=source_relative:model/ asset.proto

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(NewAssetsContact())
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}