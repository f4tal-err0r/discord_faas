package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/f4tal-err0r/discord_faas/api"
	"github.com/f4tal-err0r/discord_faas/api/context"
	cauth "github.com/f4tal-err0r/discord_faas/api/context/auth"
	"github.com/f4tal-err0r/discord_faas/api/deploy"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/db"
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)
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
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		if err := createDirIfNotExist(cfg.Filestore); err != nil {
			log.Fatalf("unable to create dir %s: %v", cfg.Filestore, err)
		}

		jwtsvc, err := security.NewJWT()
		if err != nil {
			log.Fatalf("unable to create jwt service: %v", err)
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
			cauth.NewAuthHandler(jwtsvc, dbot),
			context.NewHandler(dbot),
			deploy.NewHandler(),
		}

		r, err := api.NewRouter(handlers...)
		if err != nil {
			log.Fatalf("failed to create router: %v", err)
		}

		err = dbot.Session.Open()
		if err != nil {
			log.Fatalf("ERR: Unable to open discord session: %v", err)
		}

		log.Print("Bot Started...")

		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
			<-c

			dbot.Session.Close()
			log.Print("Bot Shutdown.")
		}()

		if err := http.ListenAndServe(":8085", r); err != nil {
			log.Fatal(err)

		}
	},
}

func createDirIfNotExist(dirPath string) error {
	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory and any necessary parents
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	} else {
		return nil
	}
	return nil
}
