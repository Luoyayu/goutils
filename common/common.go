package common

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"runtime/debug"
)

// RecoverRuntimeError :if err != nil do nothing
func RecoverRuntimeError(r interface{}, errIn error) (errOut error) {
	if errIn != nil {
		return errIn
	}
	switch r.(type) {
	case runtime.Error:
		errOut = errors.New(fmt.Sprint(r))
	case error:
		errOut = errors.New(fmt.Sprint(r))
	default:
		errOut = nil
	}
	if errOut != nil {
		errOut = errors.Wrap(errOut, string(debug.Stack()))
	}
	return
}
