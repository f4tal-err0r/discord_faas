package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	db, err = InitDB(cfg.Filestore + "/dfaas.db")
	if err != nil {
		log.Fatalf("ERR: Unable to create sqlitedb: %v", err)
	}

	dc := GetSession(cfg)
	log.Print("Bot Started...")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-c

		dc.Close()
		log.Print("Bot Shutdown.")
	}()

	log.Print("func cont")
	select {}
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
		return nil
	}
	return nil
}
