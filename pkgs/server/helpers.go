package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strconv"
)

func StrToUint(id string) uint16 {
	// Convert string to uint16
	i, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		log.Fatalf("Error converting string to uint16: %v", err)
	}
	return uint16(i)
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
