package auth

import (
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no token")
	}

	key := apiKey[len("ApiKey "):]
	if key == "" {
		return "", errors.New("no token")
	}

	return key, nil
}
