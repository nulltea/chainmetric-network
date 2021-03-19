package shared

import (
	"encoding/json"
	"errors"
	"hash/fnv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func Hash(value string) string {
	h := fnv.New32a()
	h.Write([]byte(value))
	return string(h.Sum32())
}

func ContainsString(value string, values []string) bool {
	contains := false

	if values == nil {
		return false
	}

	for i := range values {
		if values[i] == value {
			contains = true
			break
		}
	}

	return contains
}

func InvokeChaincodeParams(tx string, args ...string) [][]byte {
	var (
		res = make([][]byte, len(args) + 1)
	)

	res[0] = []byte(tx)

	for i, arg := range args {
		res[i + 1] = []byte(arg)
	}

	return res
}

func CrossChaincodeCall(ctx contractapi.TransactionContextInterface, chaincode, tx string, args ...string) ([]byte, error) {
	resp := ctx.GetStub().InvokeChaincode(chaincode, InvokeChaincodeParams(tx, args...), "")

	if resp.GetStatus() != 200 {
		return nil, errors.New(resp.GetMessage())
	}

	return resp.GetPayload(), nil
}

func MustEncode(v interface{}) string {
	data, err := json.Marshal(v); if err != nil {
		Logger.Error(err)
		return ""
	}

	return string(data)
}
