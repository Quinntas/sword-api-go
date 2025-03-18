package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

type Params struct {
	Value      string
	Salt       []byte
	Pepper     string
	Iterations int
	Length     int
}

// GenerateSalt creates a random salt of the given length.
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// EncryptValue generates an encrypted string given the value and a pepper.
// It uses default parameters for salt length, iterations, and derived key length.
func EncryptValue(value string, pepper string) (string, error) {
	salt, err := GenerateSalt(32)
	if err != nil {
		return "", err
	}
	params := &Params{
		Value:      value,
		Salt:       salt,
		Pepper:     pepper,
		Iterations: 10000,
		Length:     32,
	}
	return generateHash(params)
}

// DecryptParams parses an encrypted value and returns the parameters embedded within the hash.
func DecryptParams(hash string, pepper string) (*Params, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 5 || parts[0] != "sha256" {
		return nil, fmt.Errorf("invalid hash format")
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	iterations, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse iterations: %w", err)
	}

	length, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, fmt.Errorf("failed to parse length: %w", err)
	}

	return &Params{
		Salt:       salt,
		Iterations: iterations,
		Length:     length,
		Pepper:     pepper,
	}, nil
}

// CompareHash compares a plain value with an encrypted hash using constant-time comparison.
func CompareHash(value, hash, pepper string) (bool, error) {
	params, err := DecryptParams(hash, pepper)
	if err != nil {
		return false, err
	}
	params.Value = value

	recomputedHash, err := generateHash(params)
	if err != nil {
		return false, err
	}

	// Constant time comparison to defend against timing attacks
	return constantTimeCompare(hash, recomputedHash), nil
}

// Helper function to generate a hash using the specified parameters.
func generateHash(params *Params) (string, error) {
	hash := pbkdf2.Key([]byte(params.Pepper+params.Value), params.Salt, params.Iterations, params.Length, sha3.New256)
	saltHex := hex.EncodeToString(params.Salt)
	hashHex := hex.EncodeToString(hash)

	return fmt.Sprintf("sha256$%s$%d$%d$%s", saltHex, params.Iterations, params.Length, hashHex), nil
}

// Helper function to compare two strings in constant time.
func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
