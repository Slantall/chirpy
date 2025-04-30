package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), Subject: userID.String()})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if token.Valid != true {
		return uuid.Nil, fmt.Errorf("Token not valid.")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("could not parse claims")
	}
	userID := claims.Subject
	confID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}
	return confID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	tSplit := strings.Split(strings.TrimSpace(headers.Get("Authorization")), " ")
	if len(tSplit) != 2 {
		return "", fmt.Errorf("Issue getting Authorization token, Authorization did not contain bearer and token")
	}
	if strings.ToLower(tSplit[0]) != "bearer" {
		return "", fmt.Errorf("Issue getting Authorization token, Bearer not found")
	}
	return tSplit[1], nil
}
