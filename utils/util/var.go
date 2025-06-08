package util

import (
	"net/http"
)

//nolint:gochecknoglobals // must be global
var (
	handler http.Handler
)

func SetHandler(h http.Handler) {
	handler = h
}

func GetHandler() http.Handler {
	return handler
}
