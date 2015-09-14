package gogoapi

import (
	"log"
	"net/http"
	"time"
)

func Logger(request *http.Request, inner HandlerFunc) (int, interface{}, http.Header) {
	start := time.Now()
	status, responseObject, header := inner(request)
	log.Printf(
		" %12s  %3d  %-6s %s\n",
		time.Since(start),
		status,
		request.Method,
		request.RequestURI,
	)
	return status, responseObject, header
}
