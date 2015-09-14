package gogoapi

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*http.Request) (int, interface{}, http.Header)

type WrapperFunc func(*http.Request, HandlerFunc) (int, interface{}, http.Header)

func MethodNotAllowed(
	httpRequest *http.Request,
) (int, interface{}, http.Header) {
	status := http.StatusMethodNotAllowed
	error := fmt.Sprintf("%s method not allowed.", httpRequest.Method)
	return status, JSONError{status, error}, nil
}
