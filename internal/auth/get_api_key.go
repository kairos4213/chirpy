package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key := strings.TrimPrefix(headers.Get("authorization"), "ApiKey")
	key = strings.TrimSpace(key)
	if key == "" {
		return "", errors.New("no key found")
	}
	return key, nil
}
