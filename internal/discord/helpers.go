package discord

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strconv"
)

func strToInt(id string) int {
	// Convert string to int
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(i)
}

// GenerateRandomHash generates a random 8-character hexadecimal hash
func generateRandomHash() (string, error) {
	// 4 bytes will result in 8 hex characters
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Convert bytes to a hexadecimal string
	return hex.EncodeToString(bytes), nil
}
