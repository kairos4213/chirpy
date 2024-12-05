package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := strings.TrimPrefix(headers.Get("authorization"), "Bearer")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", errors.New("no bearer token found")
	}
	return tokenString, nil
}
