package utils

import (
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
)

// MustExecute executes `fn` function and in case of error logs it, followed by a call to os.Exit(1).
// Use `msg` to specify details of error to wrap by.
func MustExecute(fn func() error, msg string) {
	if err := fn(); err != nil {
		core.Logger.Fatal(errors.Wrap(err, msg))
	}
}

// Execute executes `fn` function and in case of error logs it.
// Use `msg` to specify details of error to wrap by.
func Execute(fn func() error, msg string) {
	if err := fn(); err != nil {
		core.Logger.Error(errors.Wrap(err, msg))
	}
}
