package jsonwebtoken

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Sign generates a signed JSON Web Token containing the provided data and expiration time.
func Sign[T any](data T, expirationTime time.Duration, secret string) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Decode parses a signed JSON Web Token and extracts the data into the provided type.
func Decode[T any](tokenString string, secret string) (*T, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Validate the token
	if !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from JWT token")
	}

	// Extract and marshal the "data" field
	data, ok := claims["data"]
	if !ok {
		return nil, errors.New("data field missing in JWT claims")
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims data: %w", err)
	}

	// Unmarshal JSON into the provided type
	var result T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims data: %w", err)
	}

	// Successful decoding
	return &result, nil
}
