package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strconv"
)

func StrToInt(id string) int64 {
	// Convert string to int
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// GenerateRandomHash generates a random 16-character hexadecimal hash
func GenerateRandomHash() (string, error) {
	// 8 bytes will result in 16 hex characters
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Convert bytes to a hexadecimal string
	return hex.EncodeToString(bytes), nil
}
