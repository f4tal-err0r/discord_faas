package main

import (
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
	"github.com/f4tal-err0r/discord_faas/pkgs/storage"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

		// Creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error creating in-cluster config: %v", err)
		}

		// Create the Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %v", err)
		}

		// create minio storage
		storage, err := storage.NewMinio()
		if err != nil {
			log.Fatalf("failed to create minio storage: %v", err)
		}

		handlers := []api.RouterAdder{
			cauth.NewAuthHandler(jwtsvc, dbot),
			context.NewHandler(dbot),
			deploy.NewHandler(cfg, dbot, clientset, storage),
		}

		r, err := api.NewRouter(jwtsvc, handlers...)
		if err != nil {
			log.Fatalf("failed to create router: %v", err)
		}

		err = dbot.Session.Open()
		if err != nil {
			log.Fatalf("ERR: Unable to open discord session: %v", err)
		}

		log.Print("Bot Started...")

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
