package client

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

func GetCurrentUser() *discordgo.User {
	token, err := NewUserAuth().GetToken()
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}
	session, err := discordgo.New("Bearer " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	user, err := session.User("@me")
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
	}

	return user
}

func FetchCacheDir(name string) string {
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
