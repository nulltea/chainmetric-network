package shared

import (
	"os"
)

var (
	// ChaincodeID stores value of "CHAINCODE_CCID" env variable.
	ChaincodeID = os.Getenv("CHAINCODE_CCID")
	// ChaincodeAddress stores value of "CHAINCODE_ADDRESS" env variable.
	ChaincodeAddress = os.Getenv("CHAINCODE_ADDRESS")
	// ChaincodeName stores value of "CHAINCODE_NAME" env variable.
	ChaincodeName = os.Getenv("CHAINCODE_NAME")
	// LogLevel stores value of "LOGGING" env variable.
	LogLevel = os.Getenv("LOGGING")
)
