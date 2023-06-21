package util

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type Response struct {
	Status            Status
	Data              string
	ErrorMessage      string
	MessageParameters []string
}

type Status string

const (
	ResponseSuccess = "success" //All went well, and (usually) some data was returned. (Required fields: status, data)
	ResponseError   = "error"   //An error occurred in processing the request, i.e. an exception was thrown. (Required: status, message	- Optionals: code, data)
)

func SetErrorResponse(w http.ResponseWriter, error error) {
	setResponse(w, Response{Status: ResponseError, ErrorMessage: error.Error()})
}

func SetAuthErrorResponse(w http.ResponseWriter, error error) {
	setAuthErrorResponse(w, error)
}

func setAuthErrorResponse(w http.ResponseWriter, error error) {
	http.Error(w, error.Error(), http.StatusForbidden)
}

func SetSuccessResponse(w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SetSuccessResponseWithoutParsingBody(w, string(data))
}

func SetSuccessResponseWithoutParsingBody(w http.ResponseWriter, body string) {
	response := Response{Status: ResponseSuccess}
	response.Data = string(body)
	setResponse(w, response)
}

func setResponse(w http.ResponseWriter, response Response) {
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	claims := jwt.MapClaims{}
	var tokenHeader = r.Header.Get("Authorization")
	tokenHeader = strings.TrimSpace(tokenHeader)

	token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	return token.Raw, err
}
