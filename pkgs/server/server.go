package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var (
	cfg      *config.Config
	db       *sql.DB
	dsession *discordgo.Session
	jwtsvc   *JWTService
)

func Start() {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatalf("ERR: Unable to create config: %v", err)
	}

	if err := createDirIfNotExist("/opt/dfaas"); err != nil {
		log.Fatalf("ERR: Unable to create dir %s: %v", cfg.Filestore, err)
	}
	db, err = NewDB(cfg)
	if err != nil {
		log.Fatalf("ERR: Unable to create db: %v", err)
	}
	err = InitDB(db)
	if err != nil {
		log.Fatalf("ERR: Unable to create sqlitedb: %v", err)
	}

	dsession, err = discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}

	if err := dsession.Open(); err != nil {
		log.Fatalf("ERR: Unable to open discord session: %v", err)
	}

	if err := RegisterCommands(db, dsession); err != nil {
		log.Printf("ERR: Unable to register commands: %v", err)
	}

	log.Print("Bot Started...")

	dsession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if _, ok := defaultCommands[i.ApplicationCommandData().Name]; ok {
				s.InteractionRespond(i.Interaction, defaultCommands[i.ApplicationCommandData().Name].Function(i.Interaction))
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		}
	})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-c

		dsession.Close()
		log.Print("Bot Shutdown.")
	}()

	jwtsvc = NewJWT()

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
