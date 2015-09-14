// Reference:
// https://gist.github.com/cryptix/45c33ecf0ae54828e63b
// http://jwt.io

package gogoapi

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	DEFAULT_TOKEN_DURATION = 24 * 60     // 24 hours, probably too long
	READ_BUFFER_SIZE       = 1024 * 1024 // 1 meg
)

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUserInfo struct {
	Username string `json:"username"`
	Type     string `json:"type"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthResource struct {
	verifyKey     *rsa.PublicKey
	signKey       *rsa.PrivateKey
	tokenDuration time.Duration
}

func NewAuthResource(privateKeyPath, publicKeyPath string, tokenDuration time.Duration) *AuthResource {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	fatal(err)

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	signBytes, err := ioutil.ReadFile(privateKeyPath)
	fatal(err)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	if tokenDuration < 1 {
		tokenDuration = DEFAULT_TOKEN_DURATION
	}
	return &AuthResource{verifyKey, signKey, tokenDuration}
}

// get token
func (auth *AuthResource) Post(request *http.Request) (int, interface{}, http.Header) {
	username := "username"
	password := "password"
	var credentials AuthCredentials
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	fatal(err)

	err = request.Body.Close()
	fatal(err)

	if err := json.Unmarshal(body, &credentials); err != nil {
		status := HTTP_UNPROCESSABLE
		return status, JSONError{status, "Unprocessable entity."}, nil
	}

	// check values
	if credentials.Username != username || credentials.Password != password {
		status := http.StatusForbidden
		return status, JSONError{status, "Authentication failed."}, nil
	}

	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims["AccessToken"] = "level 1"
	t.Claims["CustomUserInfo"] = AuthUserInfo{
		Username: credentials.Username,
		Type:     "user",
	}

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	t.Claims["exp"] = time.Now().Add(time.Minute * auth.tokenDuration).Unix()
	tokenString, err := t.SignedString(auth.signKey)
	if err != nil {
		status := http.StatusInternalServerError
		return status, JSONError{status, "Token signing error."}, nil
	}

	response := AuthResponse{
		AccessToken: tokenString,
	}

	return http.StatusOK, response, nil
}

// authorization test
// need to somehow pass AuthResource object
func (auth *AuthResource) IsAuthorized(request *http.Request) (bool, *jwt.Token, error) {
	tokenString := request.Header.Get("Authorization")
	if "" == tokenString {
		return false, nil, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.verifyKey, nil
	})

	switch err {
	case nil:
		return token.Valid, token, nil
	default:
		return false, token, err
	}
}

// canned failure response
func AuthorizationFailed() (int, interface{}, http.Header) {
	status := http.StatusUnauthorized
	return status, JSONError{status, "Access denied."}, nil
}

// authorization wrapper
func (auth *AuthResource) AuthorizationRequired(request *http.Request, inner HandlerFunc) (int, interface{}, http.Header) {
	authorized, _, _ := auth.IsAuthorized(request)
	if !authorized {
		return AuthorizationFailed()
	}
	return inner(request)
}
