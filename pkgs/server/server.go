package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var cfg *config.Config
var db *sql.DB

func Start() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ERR: Unable to create config: %v", err)
	}

	if err := createDirIfNotExist("/opt/dfaas"); err != nil {
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

	dc := GetSession(cfg)

	if err := RegisterCommands(db, dc); err != nil {
		log.Printf("ERR: Unable to register commands: %v", err)
	}

	log.Print("Bot Started...")

	dc.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if slices.ContainsFunc(defaultCommands, func(c discordgo.ApplicationCommand) bool {
				return c.Name == i.ApplicationCommandData().Name
			}) {

				if i.ApplicationCommandData().Name == "login" {
					s.InteractionRespond(i.Interaction, loginCommand(i.Interaction))
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Pong!",
					},
				})
			}
		}
	})

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
