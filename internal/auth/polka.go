package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	tSplit := strings.Split(strings.TrimSpace(headers.Get("Authorization")), " ")
	if len(tSplit) != 2 {
		return "", fmt.Errorf("Issue getting Authorization token, Authorization did not contain ApiKey and token")
	}
	if strings.ToLower(tSplit[0]) != "apikey" {
		return "", fmt.Errorf("Issue getting Authorization token, ApiKey not found")
	}
	return tSplit[1], nil
}
