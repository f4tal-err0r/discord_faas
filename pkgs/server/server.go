package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var cfg *config.Config
var db *sql.DB

func Start() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ERR: Unable to create config: %v", err)
	}

	InitServer(cfg)

	dc := GetSession(cfg)
	log.Print("Bot Started...")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-c

		dc.Close()
		log.Print("Bot Shutdown.")
	}()

	// start api server
	router := Router()
	log.Fatal(http.ListenAndServe(":8085", router))
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

func InitServer(cfg *config.Config) {
	if err := createDirIfNotExist(cfg.Filestore); err != nil {
		log.Fatalf("ERR: Unable to create dir %s: %v", cfg.Filestore, err)
	}
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("ERR: Unable to create db: %v", err)
	}
	err = InitDB(db)
	if err != nil {
		log.Fatalf("ERR: Unable to create sqlitedb: %v", err)
	}
}
