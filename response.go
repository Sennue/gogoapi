package gogoapi

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
