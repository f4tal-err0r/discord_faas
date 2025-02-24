package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/f4tal-err0r/discord_faas/pkgs/api"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/db"
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
)

type Handler struct {
	Db     *db.DBHandler
	Bot    *discord.Client
	Cfg    *config.Config
	Jwtsvc *security.JWTService
}

func NewHandler(cfg *config.Config) (*Handler, error) {
	if err := createDirIfNotExist("/opt/dfaas"); err != nil {
		return nil, fmt.Errorf("unable to create dir %s: %w", cfg.Filestore, err)
	}

	dbc, err := db.NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create db: %w", err)
	}

	if err := dbc.InitDB(); err != nil {
		return nil, fmt.Errorf("unable to create sqlitedb: %w", err)
	}
	ds, err := discord.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create discord client: %w", err)
	}

	jwtsvc, err := security.NewJWT()
	if err != nil {
		return nil, fmt.Errorf("unable to create jwt service: %w", err)
	}

	return &Handler{
		Db:     dbc,
		Bot:    ds,
		Cfg:    cfg,
		Jwtsvc: jwtsvc,
	}, nil
}

func (h *Handler) Start() {
	if err := h.Bot.Session.Open(); err != nil {
		log.Fatalf("ERR: Unable to open discord session: %v", err)
	}

	if err := h.Bot.RegisterCommands(h.Db); err != nil {
		log.Printf("ERR: Unable to register commands: %v", err)
	}

	log.Print("Bot Started...")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-c

		h.Bot.Session.Close()
		log.Print("Bot Shutdown.")
	}()

	// start api server
	apiHandler := api.NewRouter(h.Bot, h.Jwtsvc)
	log.Fatal(http.ListenAndServe(":8085", apiHandler.Router))
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
