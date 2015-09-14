// TODO: Add custom 404 so missing pages can be logged.
// Reference:
// http://stackoverflow.com/questions/26141953/custom-404-with-gorilla-mux-and-std-http-fileserver

package gogoapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

type HandlerFunc func(*http.Request) (int, interface{}, http.Header)

type WrapperFunc func(*http.Request, HandlerFunc) (int, interface{}, http.Header)

type API struct {
	router           *mux.Router
	wrapperFuncs     []WrapperFunc
	methodNotAllowed HandlerFunc
}

func NewAPI(wrapperFuncs []WrapperFunc) *API {
	router := mux.NewRouter().StrictSlash(true)
	api := API{router, wrapperFuncs, nil}
	api.methodNotAllowed = api.WrapHandler(MethodNotAllowed, wrapperFuncs)
	return &api
}

func MethodNotAllowed(
	httpRequest *http.Request,
) (int, interface{}, http.Header) {
	status := http.StatusMethodNotAllowed
	error := fmt.Sprintf("%s method not allowed.", httpRequest.Method)
	return status, JSONError{status, error}, nil
}

func (api *API) HttpRequestHandler(
	inner HandlerFunc,
) http.HandlerFunc {
	if nil == inner {
		inner = MethodNotAllowed
	}
	handler := func(
		httpResponse http.ResponseWriter,
		httpRequest *http.Request,
	) {
		status, responseObject, header := inner(httpRequest)
		if _, supported := responseObject.(StatusResponse); !supported {
			responseObject = JSONData{status, responseObject}
		}
		jsonResponse, err := json.MarshalIndent(responseObject, "", "  ")
		if nil != err {
			header = nil
			status = http.StatusInternalServerError
			message := fmt.Sprintf(
				"{\n  \"status\": %d,\n  \"error\": \"Internal server error. %s\"\n}",
				status,
				err.Error(),
			)
			jsonResponse = []byte(message)
		}
		for key, values := range header {
			for _, value := range values {
				httpResponse.Header().Add(key, value)
			}
		}
		httpResponse.Header().Set(
			"Content-Type", "application/json; charset=UTF-8",
		)
		httpResponse.WriteHeader(status)
		httpResponse.Write(jsonResponse)
		httpResponse.Write([]byte("\n"))
	}
	return handler
}

func (api *API) AddRoute(
	method, path string,
	inner HandlerFunc,
) {
	handler := api.HttpRequestHandler(inner)
	api.router.Methods(method).Path(path).Handler(handler)
}

func (api *API) WrapHandler(
	handler HandlerFunc,
	wrapperFuncs []WrapperFunc,
) HandlerFunc {
	if nil == wrapperFuncs || 0 == len(wrapperFuncs) {
		return handler
	}
	// every layer is a HandlerFunc that calls a WrapperFunc
	for _, nextWrapperFunc := range wrapperFuncs {
		wrapperFunc := nextWrapperFunc
		innerHandler := handler
		handler = func(
			httpRequest *http.Request,
		) (int, interface{}, http.Header) {
			return wrapperFunc(httpRequest, innerHandler)
		}
	}
	return handler
}

func (api *API) AddResource(
	resource interface{},
	path string,
	resourceWrapperFuncs []WrapperFunc,
) {
	wrapperFuncs := append(resourceWrapperFuncs, api.wrapperFuncs...)
	var handler HandlerFunc

	if resource, supported := resource.(GetResource); supported {
		handler = api.WrapHandler(resource.Get, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(GET, path, handler)

	if resource, supported := resource.(PostResource); supported {
		handler = api.WrapHandler(resource.Post, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(POST, path, handler)

	if resource, supported := resource.(PutResource); supported {
		handler = api.WrapHandler(resource.Put, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(PUT, path, handler)

	if resource, supported := resource.(DeleteResource); supported {
		handler = api.WrapHandler(resource.Delete, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(DELETE, path, handler)

	if resource, supported := resource.(HeadResource); supported {
		handler = api.WrapHandler(resource.Head, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(HEAD, path, handler)

	if resource, supported := resource.(PatchResource); supported {
		handler = api.WrapHandler(resource.Patch, wrapperFuncs)
	} else {
		handler = api.methodNotAllowed
	}
	api.AddRoute(PATCH, path, handler)
}

func (api *API) Start(host string, port int) error {
	networkAddress := fmt.Sprintf("%s:%d", host, port)
	return http.ListenAndServe(networkAddress, api.router)
}
