package shared

import (
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// BootstrapContract performs start sequence of the Smart Contract handler.
func BootstrapContract(contract contractapi.ContractInterface) {
	chaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		Logger.Fatalf("Error creating new Smart Contract: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID: ChaincodeID,
		Address: ChaincodeAddress,
		CC: chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	Logger.Info("Contract is up and running...")

	if err = server.Start(); err != nil {
		Logger.Fatalf("Error starting %s chaincode: %s", os.Getenv("CHAINCODE_NAME"), err)
	}
}
