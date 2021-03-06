package shared

import (
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func BootstrapContract(contract contractapi.ContractInterface) {
	chaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		Logger.Fatalf("Error creating new Smart Contract: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID: os.Getenv("CHAINCODE_CCID"),
		Address: os.Getenv("CHAINCODE_ADDRESS"),
		CC: chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err = server.Start(); err != nil {
		Logger.Fatalf("Error starting %s chaincode: %s", os.Getenv("CHAINCODE_NAME"), err)
	}

	Logger.Info("Contract up and running...")
}
