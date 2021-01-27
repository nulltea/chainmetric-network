package main

//go:generate protoc --proto_path=model/ --gogofaster_out=paths=source_relative:model/ asset.proto

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(NewAssetsContact())
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID: os.Getenv("CHAINCODE_CCID"),
		Address: os.Getenv("CHAINCODE_ADDRESS"),
		CC: chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	err = server.Start()
	if err != nil {
		fmt.Printf("Error starting %s chaincode: %s", os.Getenv("CHAINCODE_NAME"), err)
	}
}