package server_test

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ewohltman/discordgo-mock/mockrest"
	"github.com/ewohltman/discordgo-mock/mocksession"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/server"
)

var db *sql.DB
var session *discordgo.Session

func init() {
	var err error
	// Use the config struct instead of config.New() and set datastore to memory
	cnfg := &config.Config{
		DBPath: ":memory:",
	}

	db, err = server.NewDB(cnfg)
	if err != nil {
		log.Fatalf("failed to create db: %v", err)
	}
	err = server.InitDB(db)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// Mock a new discord session and create a new guild
	state, err := newState()
	if err != nil {
		log.Fatal(err)
	}
	session, err = mocksession.New(
		mocksession.WithState(state),
		mocksession.WithClient(&http.Client{
			Transport: mockrest.NewTransport(state),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

}

func TestInitDB(t *testing.T) {
	// Test that the GuildMetadata, ApprovedRoles, and Commands tables exist
	_, err := db.Exec("SELECT * FROM GuildMetadata")
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

func TestCmdsDb(t *testing.T) {
	hash, err := server.GenerateRandomHash()
	guildid := generateSnowflake()

	if err != nil {
		t.Fatal(err)
	}

	time := time.Now()

	commandRow := []server.CommandsTableRow{
		{
			Command:       "test",
			Hash:          hash,
			Last_modified: &time,
			Guildid:       guildid,
		},
		{
			Command:       "test2",
			Hash:          hash,
			Last_modified: &time,
			Guildid:       guildid,
		},
		{
			Command:       "test3",
			Hash:          hash,
			Last_modified: &time,
			Guildid:       generateSnowflake(),
		},
	}

	// Use RowWriter to insert a new command, generate a random md5 hash for the hash field

	for _, row := range commandRow {
		err = server.RowWriter(db, row)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err != nil {
		t.Fatal(err)
	}

	// Call GetCmdsDb
	cmds, err := server.GetCmdsDb(db, guildid)
	if err != nil {
		t.Fatal(err)
	}

	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cmds))
	}

	if cmds[0].Command != "test" {
		t.Fatalf("expected 'test', got %s", cmds[0].Command)
	}

	if cmds[1].Command != "test2" {
		t.Fatalf("expected 'test2', got %s", cmds[1].Command)
	}

}

func generateSnowflake() int64 {
	var DiscordEpoch int64 = 1420070400000
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := currentTime - DiscordEpoch

	// Generate a random 10-bit machine ID and a random 12-bit sequence number
	var randomBytes [2]byte
	_, err := rand.Read(randomBytes[:])
	if err != nil {
		panic(err)
	}
	machineID := int64(randomBytes[0]) & 0x03FF // Mask to 10 bits
	sequence := int64(randomBytes[1]) & 0x0FFF  // Mask to 12 bits

	// Combine the parts into a 64-bit integer
	snowflake := (timestamp << 22) | (machineID << 12) | sequence

	return snowflake
}

func TestGetRolebyGuildid(t *testing.T) {
	// Create a in memory database
	db, err := server.NewDB(&config.Config{
		DBPath: ":memory:",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = server.InitDB(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new guild
	guild := server.GuildMetaRow{
		Guildid:  generateSnowflake(),
		Owner:    "test",
		Name:     "test",
		Textchan: "test",
		Source:   "local",
	}
	err = server.RowWriter(db, guild)
	if err != nil {
		t.Fatal(err)
	}

	// Create multiple new roles and write to db
	role := []server.ApprovedRolesRow{
		{
			Roleid:  "1",
			Guildid: guild.Guildid,
		},
		{
			Roleid:  "2",
			Guildid: guild.Guildid,
		},
		{
			Roleid:  "3",
			Guildid: guild.Guildid,
		},
		{
			Roleid:  "4",
			Guildid: generateSnowflake(),
		},
	}

	for _, v := range role {
		err = server.RowWriter(db, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Call GetRolebyGuildid and check if it returns the correct roles
	roles, err := server.GetRolebyGuildid(db, guild.Guildid)
	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(roles))
	}

	if roles[0] != "1" || roles[1] != "2" || roles[2] != "3" {
		t.Errorf("Expected roles to be 1, 2, and 3, got %v", roles)
	}

	// Check if GetRolebyGuildid returns an error if guildid is not found
	_, err = server.GetRolebyGuildid(db, int64(3))
	if err == nil {
		t.Error("Expected error, got nil")
	}

}
