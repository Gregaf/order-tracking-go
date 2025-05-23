package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	charSet = letters + digits
)

var (
	csLen    = big.NewInt(int64(len(charSet)))
	digitLen = big.NewInt(int64(len(digits)))
)

func GenerateRandomSegment(length int) (string, error) {
	if length < 2 {
		return "", fmt.Errorf("segment length must be at least 2 characters")
	}

	result := make([]byte, length)

	randomIndex, err := rand.Int(rand.Reader, digitLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate random index: %w", err)
	}

	digitPosition, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random position: %w", err)
	}

	result[digitPosition.Int64()] = digits[randomIndex.Int64()]

	for i := 0; i < length; i++ {
		if i == int(digitPosition.Int64()) {
			continue
		}

		randomIndex, err := rand.Int(rand.Reader, csLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}
		result[i] = charSet[randomIndex.Int64()]
	}

	return string(result), nil
}

func IsOneOf[T comparable](a T, vals ...T) bool {
	for _, s := range vals {
		if s == a {
			return true
		}
	}
	return false
}
