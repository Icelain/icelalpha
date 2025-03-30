package jwtauth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateJWTToken(email string, secretKey string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,                            // Subject (user identifier)
		"iss": "icealpha",                       // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
