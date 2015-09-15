package gogoapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	READ_BUFFER_SIZE = 1024 * 1024 // 1 meg
)

type DummyAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DummyAuthUserInfo struct {
	Username string `json:"username"`
	Type     string `json:"type"`
}

func DummyAuthValidator(request *http.Request) (bool, map[string]interface{}, StatusResponse) {
	username := "username"
	password := "password"

	var credentials DummyAuthCredentials
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	fatal(err)

	err = request.Body.Close()
	fatal(err)

	if err := json.Unmarshal(body, &credentials); err != nil {
		return false, nil, JSONError{HTTP_UNPROCESSABLE, "Unprocessable entity."}
	}

	// check values
	if credentials.Username != username || credentials.Password != password {
		return false, nil, JSONError{http.StatusForbidden, "Authentication failed."}
	} else {
		claims := make(map[string]interface{})
		claims["AccessToken"] = "level 1"
		claims["CustomUserInfo"] = DummyAuthUserInfo{
			Username: credentials.Username,
			Type:     "user",
		}
		return true, claims, nil
	}
}
