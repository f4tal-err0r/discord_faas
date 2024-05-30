package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var db *sql.DB

func Start() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ERR: Unable to create config: %v", err)
	}
	if err := createDirIfNotExist(cfg.Filestore); err != nil {
		log.Fatalf("ERR: Unable to create dir %s: %v", cfg.Filestore, err)
	}
	_, err = InitDB(cfg.Filestore)
	if err != nil {
		log.Fatalf("ERR: Unable to create sqlitedb: %v", err)
	}

}

func createDirIfNotExist(dirPath string) error {
	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory and any necessary parents
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		fmt.Printf("Directory %s created successfully\n", dirPath)
	} else {
		fmt.Printf("Directory %s already exists\n", dirPath)
	}
	return nil
}
