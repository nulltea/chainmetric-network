package core

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/spf13/viper"
)

// BootstrapChaincodeServer performs start sequence of the Smart Contract handler.
func BootstrapChaincodeServer(contract contractapi.ContractInterface) {
	chaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		Logger.Fatalf("Error creating new Smart Contract: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID: viper.GetString("chaincode.ccid"),
		Address: viper.GetString("chaincode.address"),
		CC: chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	Logger.Info("Contract is up and running...")

	if err = server.Start(); err != nil {
		Logger.Fatalf("Error starting %s chaincode: %s", viper.GetString("chaincode.name"), err)
	}
}
