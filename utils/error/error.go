// Package apierror provides common error definitions for API testing operations.
package apierror

import "errors"

var (
	// ErrIndexOutOfRange is returned when an index is out of range.
	ErrIndexOutOfRange = errors.New("index out of range")
	// ErrIndexNotFound is returned when an index is not found.
	ErrIndexNotFound = errors.New("index not found")
	// ErrIndexInvaild is returned when an index is not a valid number.
	ErrIndexInvaild = errors.New("index is not a number")
	// ErrTypeNotSupported is returned when a type is not supported.
	ErrTypeNotSupported = errors.New("type not supported")

	// ErrDataFieldNotFound is returned when a data field is not found.
	ErrDataFieldNotFound = errors.New("data field not found")
	// ErrDataTypeMustBeMap is returned when data type must be a map.
	ErrDataTypeMustBeMap = errors.New("data type must be map")

	// ErrStepRespNotFound is returned when a step response is not found.
	ErrStepRespNotFound = errors.New("step response not found")

	// ErrVarsNotFound is returned when variables are not found.
	ErrVarsNotFound = errors.New("vars not found")
)
