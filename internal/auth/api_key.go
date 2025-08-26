package auth

import (
	"net/http"
	"strings"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	splitHeader := strings.Split(authHeader, "ApiKey ")
	if len(splitHeader) < 2 {
		return "", errors.New("No ApiKey present")
	}
	return splitHeader[1], nil
}
