package utils

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
)

// Iterate performs iteration over state query results calling `fn` for each record.
func Iterate(
	iter shim.StateQueryIteratorInterface,
	fn func(key string, value []byte) error,
) {
	defer func() {
		if err := iter.Close(); err != nil {
			core.Logger.Error(errors.Wrap(err, "failed to close state query iterator"))
		}
	}()

	for iter.HasNext() {
		result, err := iter.Next(); if err != nil {
			core.Logger.Error(errors.Wrap(err, "failed to iterate over records"))
			continue
		}

		if err := fn(result.Key, result.Value); err != nil {
			core.Logger.Error(err)
		}
	}
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
	defer core.Logger.Error(err)

	return err
}

// LoggedErrorf wraps `err` error with formatted message and logs it.
func LoggedErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	err = errors.Wrapf(err, format, args)
	defer core.Logger.Error(err)

	return err
}
