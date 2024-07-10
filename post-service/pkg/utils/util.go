package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractSubFromToken(tokenString string) (string, error) {
	// Parse the token without verifying the signature
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

	if err != nil {
		return "", errors.New(fmt.Sprintf("Error parsing token: %v", err))
	}

	// Access the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", errors.New("subject (sub) claim not found in token")
	}

	return "", errors.New("failed to parse claims")
}
