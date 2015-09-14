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

func PageNotFound(
	httpRequest *http.Request,
) (int, interface{}, http.Header) {
	status := http.StatusNotFound
	host := httpRequest.URL.Host
	path := httpRequest.URL.Path
	error := fmt.Sprintf("%s%s page not found.", host, path)
	return status, JSONError{status, error}, nil
}

