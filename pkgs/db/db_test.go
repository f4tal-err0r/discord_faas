package db_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/f4tal-err0r/discord_faas/pkgs/db"
)

var (
	testDB *DBHandler
)

func init() {
	db, err := NewDB(":memory:")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}
	testDB = db
}

func TestInsertGuildMetadata(t *testing.T) {
	guildID := rand.Int()
	guild := GuildMetadata{GuildID: guildID, Source: "test", Name: "test", Owner: fmt.Sprintf("foo-%d", guildID)}
	if err := testDB.InsertGuildMetadata(guild); err != nil {
		t.Fatalf("InsertGuildMetadata: %v", err)
	}

	retrievedGuild, err := testDB.GetGuild(guildID)
	if err != nil {
		t.Fatalf("GetGuild: %v", err)
	}

	if *retrievedGuild != guild {
		t.Fatalf("expected %v, got %v", guild, retrievedGuild)
	}
}

func TestInsertFunction(t *testing.T) {
	guildID := rand.Int()
	function := Function{Name: fmt.Sprintf("test-%d", guildID), Description: "test", Runtime: "golang", GuildID: guildID}
	if err := testDB.InsertFunction(function); err != nil {
		t.Fatalf("InsertFunction: %v", err)
	}

	functions, err := testDB.GetFunctionsByGuild(guildID)
	if err != nil {
		t.Fatalf("GetFunctionsByGuild: %v", err)
	}

	if len(functions) == 0 || functions[0].Name != function.Name {
		t.Fatalf("expected function %v, got %v", function, functions)
	}
}

func TestInsertCommand(t *testing.T) {
	guildID := rand.Int()
	command := Command{Command: fmt.Sprintf("test-%d", guildID), Description: "test", GuildID: guildID}
	if err := testDB.InsertCommand(command); err != nil {
		t.Fatalf("InsertCommand: %v", err)
	}

	commands, err := testDB.GetCommandsByGuild(guildID)
	if err != nil {
		t.Fatalf("GetCommandsByGuild: %v", err)
	}

	if len(commands) == 0 || commands[0].Command != command.Command {
		t.Fatalf("expected command %v, got %v", command, commands)
	}
}

func TestInsertCommandArgument(t *testing.T) {
	commandID := rand.Int()
	arg := CommandArgument{CommandID: commandID, Name: fmt.Sprintf("test-%d", commandID), Description: "test"}
	if err := testDB.InsertCommandArgument(arg); err != nil {
		t.Fatalf("InsertCommandArgument: %v", err)
	}

	args, err := testDB.GetArgumentsByCommand(commandID)
	if err != nil {
		t.Fatalf("GetArgumentsByCommand: %v", err)
	}

	if len(args) == 0 || args[0].Name != arg.Name {
		t.Fatalf("expected argument %v, got %v", arg, args)
	}
}
