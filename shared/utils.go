package shared

import (
	"hash/fnv"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

// Hash returns the hash-sum of the given object.
// It uses FNV32.
func Hash(value string) string {
	h := fnv.New32a()
	h.Write([]byte(value))
	return strconv.Itoa(int(h.Sum32()))
}

// InvokeChaincodeParams builds chaincode invoke params object.
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

// CrossChaincodeCall performs cross-chaincode call.
func CrossChaincodeCall(ctx contractapi.TransactionContextInterface, chaincode, tx string, args ...string) ([]byte, error) {
	resp := ctx.GetStub().InvokeChaincode(chaincode, InvokeChaincodeParams(tx, args...), "")

	if resp.GetStatus() != 200 {
		return nil, errors.New(resp.GetMessage())
	}

	return resp.GetPayload(), nil
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
