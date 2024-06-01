package server_test

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
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
	state, err := newState()
	if err != nil {
		t.Fatal(err)
	}

	hash, err := server.GenerateRandomHash()
	guildid := GenerateRandomUint16()

	if err != nil {
		t.Fatal(err)
	}

	time := time.Now()

	commandRow := server.CommandsTableRow{
		Command:       "test",
		Hash:          hash,
		Guildid:       guildid,
		Last_modified: &time,
		Desc:          "long description for test command",
	}

	// Use RowWriter to insert a new command, generate a random md5 hash for the hash field
	err = server.RowWriter(db, commandRow)

	if err != nil {
		t.Fatal(err)
	}

	// Call GetCmdsDb
	cmds, err := server.GetCmdsDb(db, "test", guildid)
	if err != nil {
		t.Fatal(err)
	}

	if cmds.Command != "test" {
		t.Errorf("Expected command to be %s, got %s", "test", cmds.Command)
	}

	if cmds.Hash != hash {
		t.Errorf("Expected hash to be %s, got %s", hash, cmds.Hash)
	}

	if cmds.Guildid != guildid {
		t.Errorf("Expected guildid to be %v, got %v", state.Guilds[0].ID, cmds.Guildid)
	}

	if cmds.Last_modified == nil {
		t.Error("Expected Last_modified to be not nil")
	}

	if cmds.Desc != "long description for test command" {
		t.Errorf("Expected desc to be %s, got %s", "long description for test command", cmds.Desc)
	}
}

func GenerateRandomUint16() uint16 {
	var num uint16
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return 0
	}
	num = binary.LittleEndian.Uint16(bytes)
	return num
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
		Guildid:  GenerateRandomUint16(),
		Owner:    "test",
		Name:     "test",
		Textchan: "test",
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
			Guildid: uint16(5),
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
	_, err = server.GetRolebyGuildid(db, uint16(3))
	if err == nil {
		t.Error("Expected error, got nil")
	}

}
