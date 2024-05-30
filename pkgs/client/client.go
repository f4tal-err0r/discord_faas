package client

import (
	"log"
	"os"
	"path/filepath"
)

func FetchCache(name string) string {
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Unable to fetch cache directory: %v", err)
	}
	cacheDir := cache + "/dfaas"
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		log.Fatalf("Unable to create cache directory: %v", err)
	}
	return filepath.Join(cacheDir, name+".json")
}
