package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/api"
	"github.com/f4tal-err0r/discord_faas/api/context"
	cauth "github.com/f4tal-err0r/discord_faas/api/context/auth"
	"github.com/f4tal-err0r/discord_faas/api/functions"
	"github.com/f4tal-err0r/discord_faas/internal/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/db"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	mode    string
	backend string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)

	serverCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "/app/config/config.yaml", "Path to config file")
	serverCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "kubernetes", "Mode to run in (kubernetes, docker)")
	serverCmd.PersistentFlags().StringVar(&backend, "backend", "local://data", "Storage backend to use. Supports local://path or s3://bucket")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Functions-as-a-Service Server",
	Long:  "Discord FaaS kubernetes controller",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.New(cfgPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		jwtsvc, err := security.NewJWT()
		if err != nil {
			log.Fatalf("unable to create jwt service: %v", err)
		}

		//create dir if not exists
		if _, err := os.Stat(cfg.Filestore); os.IsNotExist(err) {
			err := os.MkdirAll(cfg.Filestore, os.ModePerm)
			if err != nil {
				log.Fatalf("unable to create filestore dir: %v", err)
			}
		}

		//create dbpath dir if not exists
		if _, err := os.Stat(cfg.DBPath); os.IsNotExist(err) {
			err := os.MkdirAll(cfg.DBPath, os.ModePerm)
			if err != nil {
				log.Fatalf("unable to create dbpath dir: %v", err)
			}
		}

		dbc, err := db.NewDB(cfg.DBPath)
		if err != nil {
			log.Fatalf("unable to create db: %v", err)
		}

		dbot, err := discord.NewClient(dbc, cfg)
		if err != nil {
			log.Fatalf("failed to create discord bot: %v", err)
		}

		handlers := []api.RouterAdder{
			cauth.NewAuthHandler(jwtsvc, dbot, cfg),
			context.NewHandler(dbot, cfg),
			functions.NewHandler(cfg),
		}

		r, err := api.NewRouter(jwtsvc, handlers...)
		if err != nil {
			log.Fatalf("failed to create router: %v", err)
		}

		if err := dbot.StartBotHandler(); err != nil {
			log.Fatalf("error starting bot handler: %v", err)
		}

		// if err := purgeGlobalCommands(dbot.Session); err != nil {
		// 	log.Fatalf("error purging global commands: %v", err)
		// }

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		appserv := &http.Server{Addr: ":8080", Handler: r}

		go func() {
			if err := appserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe: %v", err)
			}
		}()

		<-stopChan

		if err := appserv.Close(); err != nil {
			log.Printf("Error closing appserver: %v", err)
		}

		dbot.Session.Close()
		log.Print("Bot Shutdown.")
	},
}

func purgeGlobalCommands(s *discordgo.Session) error {
	// Get all global commands (pass empty string for guildID)
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("failed to get global commands: %w", err)
	}

	// Delete each command
	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Printf("Failed to delete global command %s: %v", cmd.Name, err)
			continue
		}
		log.Printf("Deleted global command: %s", cmd.Name)
	}

	//delete commands per guild
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("failed to get user guilds: %w", err)
	}

	for _, guild := range guilds {
		commands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			log.Printf("Failed to get commands for guild %s: %v", guild.ID, err)
			continue
		}

		for _, cmd := range commands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, cmd.ID)
			if err != nil {
				log.Printf("Failed to delete command %s in guild %s: %v", cmd.Name, guild.ID, err)
				continue
			}
			log.Printf("Deleted command %s in guild %s", cmd.Name, guild.ID)
		}
	}
	log.Printf("Purged %d global commands", len(commands))
	return nil
}
