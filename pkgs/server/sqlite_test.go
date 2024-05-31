package server_test

import (
	"testing"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/server"
)

func TestInitDB(t *testing.T) {
	// Use the config struct instead of config.New() and set datastore to memory
	cnfg := &config.Config{
		DBPath: ":memory:",
	}

	db, err := server.NewDB(cnfg)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	err = server.InitDB(db)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	// Test that the GuildMetadata, ApprovedRoles, and Commands tables exist
	_, err = db.Exec("SELECT * FROM GuildMetadata")
	if err != nil {
		t.Fatalf("failed to select from GuildMetadata: %v", err)
	}
	_, err = db.Exec("SELECT * FROM ApprovedRoles")
	if err != nil {
		t.Fatalf("failed to select from ApprovedRoles: %v", err)
	}
	_, err = db.Exec("SELECT * FROM Commands")
	if err != nil {
		t.Fatalf("failed to select from Commands: %v", err)
	}
}
