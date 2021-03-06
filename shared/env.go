package shared

import "os"

var (
	ChaincodeID = os.Getenv("CHAINCODE_CCID")
	ChaincodeAddress = os.Getenv("CHAINCODE_ADDRESS")
	ChaincodeName = os.Getenv("CHAINCODE_NAME")
)
