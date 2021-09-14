package core

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Chaincode struct {
	*contractapi.ContractChaincode
}

func NewChaincode(contracts ...contractapi.ContractInterface) (*Chaincode, error) {
	cc, err := contractapi.NewChaincode(contracts...)
	if err != nil {
		Logger.Fatalf("Error creating new Smart Contract: %s", err)
	}

	return &Chaincode{
		ContractChaincode: cc,
	}, nil
}


func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// TODO: call init
	return cc.ContractChaincode.Init(stub)
}

func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.ContractChaincode.Invoke(stub)
}
