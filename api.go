// TODO: Add SSL support
// Reference:
// https://gist.github.com/michaljemala/d6f4e01c4834bf47a9c4

package gogoapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	router           *mux.Router
	wrapperFuncs     []WrapperFunc
	methodNotAllowed HandlerFunc
}

func NewAPI(wrapperFuncs []WrapperFunc) *API {
	router := mux.NewRouter().StrictSlash(true)
	api := API{router, wrapperFuncs, nil}
	api.methodNotAllowed = api.WrapHandler(MethodNotAllowed, wrapperFuncs)
	wrappedPageNotFound := api.WrapHandler(PageNotFound, wrapperFuncs)
	router.NotFoundHandler = api.HttpRequestHandler(wrappedPageNotFound)
	return &api
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
