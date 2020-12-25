package authentication

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// IssueJWTWithSecret issues and sign a JWT with a secret
func IssueJWTWithSecret(secret, email string, expires time.Time) (string, error) {
	key := []byte(secret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expires.Unix(),
		Subject:   email,
		Issuer:    "secProxy",
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

// ValidateJWTWithSecret checks JWT signing algorithm as well the signature
func ValidateJWTWithSecret(secret, tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return err == nil && token != nil && token.Valid
}
