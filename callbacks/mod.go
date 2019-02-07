package callbacks

import (
	"github.com/pkg/errors"
)

var DeregisterCallback = errors.New("callback deregistered")

type callback func(register registerCallback, params ...interface{}) error
type reduceCallback func(register registerReduceCallback, in interface{}, params ...interface{}) (interface{}, error)

type registerCallback func(c callback)
type registerReduceCallback func(c reduceCallback)