// TODO: Add custom 404 so missing pages can be logged.
// Reference:
// http://stackoverflow.com/questions/26141953/custom-404-with-gorilla-mux-and-std-http-fileserver

package gogoapi

import (
)

type StatusResponse interface {
	Status() int
}

type JSONError struct {
	StatusCode int    `json:"status"`
	Error      string `json:"error"`
}

type JSONMessage struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

type JSONData struct {
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data"`
}

func (response JSONError) Status() int {
	return response.StatusCode
}

func (response JSONMessage) Status() int {
	return response.StatusCode
}

func (response JSONData) Status() int {
	return response.StatusCode
}

