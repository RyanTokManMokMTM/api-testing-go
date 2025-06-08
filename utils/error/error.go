package apierror

import "errors"

var (
	ErrIndexOutOfRange  = errors.New("index out of range")
	ErrIndexNotFound    = errors.New("index not found")
	ErrIndexInvaild     = errors.New("index is not a number")
	ErrTypeNotSupported = errors.New("type not supported")

	ErrDataFieldNotFound = errors.New("data field not found")
	ErrDataTypeMustBeMap = errors.New("data type must be map")

	ErrStepRespNotFound = errors.New("step response not found")

	ErrVarsNotFound = errors.New("vars not found")
)
