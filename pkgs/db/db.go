package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// GuildMetadata represents a row in the GuildMetadata table
type GuildMetadata struct {
	GuildID int    `json:"guildid"`
	Source  string `json:"source"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
}

// Function represents a row in the Functions table
type Function struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Runtime     string    `json:"runtime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GuildID     int       `json:"guildid"`
}

// ApprovedRole represents a row in the ApprovedRoles table
type ApprovedRole struct {
	RoleID    string    `json:"roleid"`
	GuildID   int       `json:"guildid"`
	CreatedAt time.Time `json:"created_at"`
}

// Command represents a row in the Commands table
type Command struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	Description string `json:"description"`
	GuildID     int    `json:"guildid"`
}

// CommandArgument represents a row in the CommandArguments table
type CommandArgument struct {
	ID          int    `json:"id"`
	CommandID   int    `json:"command_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DBHandler struct {
	db *sql.DB
}

func NewDB(DBPath string) (*DBHandler, error) {
	var handler DBHandler

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	handler.db = db

	if err := applyMigration(&handler); err != nil {
		return nil, fmt.Errorf("failed to apply migration: %v", err)
	}

	return &handler, nil
}

func (h *DBHandler) Close() error {
	return h.db.Close()
}

// InsertGuildMetadata inserts a new guild
func (h *DBHandler) InsertGuildMetadata(guild GuildMetadata) error {
	_, err := h.db.Exec("INSERT INTO GuildMetadata (guildid, source, name, owner) VALUES (?, ?, ?, ?)",
		guild.GuildID, guild.Source, guild.Name, guild.Owner)
	return err
}

// GetGuild retrieves a guild by ID
func (h *DBHandler) GetGuild(guildID int) (*GuildMetadata, error) {
	row := h.db.QueryRow("SELECT guildid, source, name, owner FROM GuildMetadata WHERE guildid = ?", guildID)
	var guild GuildMetadata
	if err := row.Scan(&guild.GuildID, &guild.Source, &guild.Name, &guild.Owner); err != nil {
		return nil, err
	}
	return &guild, nil
}

// InsertFunction inserts a new function
func (h *DBHandler) InsertFunction(f Function) error {
	_, err := h.db.Exec("INSERT INTO Functions (name, description, runtime, guildid) VALUES (?, ?, ?, ?)",
		f.Name, f.Description, f.Runtime, f.GuildID)
	return err
}

// GetFunctionsByGuild retrieves functions for a given guild
func (h *DBHandler) GetFunctionsByGuild(guildID int) ([]Function, error) {
	rows, err := h.db.Query("SELECT id, name, description, runtime, created_at, updated_at, guildid FROM Functions WHERE guildid = ?", guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var functions []Function
	for rows.Next() {
		var f Function
		if err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.Runtime, &f.CreatedAt, &f.UpdatedAt, &f.GuildID); err != nil {
			return nil, err
		}
		functions = append(functions, f)
	}
	return functions, nil
}

// InsertCommand inserts a new command
func (h *DBHandler) InsertCommand(c Command) error {
	_, err := h.db.Exec("INSERT INTO Commands (command, description, guildid) VALUES (?, ?, ?)",
		c.Command, c.Description, c.GuildID)
	return err
}

// GetCommandsByGuild retrieves commands for a guild
func (h *DBHandler) GetCommandsByGuild(guildID int) ([]Command, error) {
	rows, err := h.db.Query("SELECT id, command, description, guildid FROM Commands WHERE guildid = ?", guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.Command, &c.Description, &c.GuildID); err != nil {
			return nil, err
		}
		commands = append(commands, c)
	}
	return commands, nil
}

// InsertCommandArgument inserts a new command argument
func (h *DBHandler) InsertCommandArgument(arg CommandArgument) error {
	_, err := h.db.Exec("INSERT INTO CommandArguments (command_id, argument, description) VALUES (?, ?, ?)",
		arg.CommandID, arg.Name, arg.Description)
	return err
}

// GetArgumentsByCommand retrieves arguments for a command
func (h *DBHandler) GetArgumentsByCommand(commandID int) ([]CommandArgument, error) {
	rows, err := h.db.Query("SELECT id, command_id, argument, description FROM CommandArguments WHERE command_id = ?", commandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var args []CommandArgument
	for rows.Next() {
		var arg CommandArgument
		if err := rows.Scan(&arg.ID, &arg.CommandID, &arg.Name, &arg.Description); err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}
