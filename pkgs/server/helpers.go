package server

import (
	"log"
	"strconv"
)

func strToUint(id string) uint64 {
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
