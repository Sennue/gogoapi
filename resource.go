package gogoapi

import (
	"net/http"
)

type GetResource interface {
	Get(*http.Request) (int, interface{}, http.Header)
}

type PostResource interface {
	Post(*http.Request) (int, interface{}, http.Header)
}

type PutResource interface {
	Put(*http.Request) (int, interface{}, http.Header)
}

type DeleteResource interface {
	Delete(*http.Request) (int, interface{}, http.Header)
}

type HeadResource interface {
	Head(*http.Request) (int, interface{}, http.Header)
}

type PatchResource interface {
	Patch(*http.Request) (int, interface{}, http.Header)
}
