// Package util provides utility functions and variables for common operations.
package util

import (
	"net/http"
)

//nolint:gochecknoglobals // must be global
var (
	handler http.Handler
)

// SetHandler sets the global HTTP handler for testing purposes.
func SetHandler(h http.Handler) {
	handler = h
}

// GetHandler returns the global HTTP handler.
func GetHandler() http.Handler {
	return handler
}
