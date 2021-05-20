package shared

import (
	"encoding/json"
	"hash/fnv"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

func Hash(value string) string {
	h := fnv.New32a()
	h.Write([]byte(value))
	return strconv.Itoa(int(h.Sum32()))
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

// LoggedError wraps `err` error with `msg` message and logs it.
func LoggedError(err error, msg string) error {
	if err == nil {
		return nil
	}

	err = errors.Wrap(err, msg)
	defer Logger.Error(err)

	return err
}

// LoggedErrorf wraps `err` error with formatted message and logs it.
func LoggedErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	err = errors.Wrapf(err, format, args)
	defer Logger.Error(err)

	return err
}
