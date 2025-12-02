package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
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

// Runtime represents a row in the Runtimes table
type Runtime struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Lang      string    `json:"lang"`
	Repo      string    `json:"repo"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	CommitID  string    `json:"commitid"`
}

// Function represents a row in the Functions table
type Function struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Hash        string    `json:"hash"`
	Description string    `json:"description"`
	Runtime     string    `json:"runtime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GuildID     int       `json:"guildid"`
	RuntimeID   int       `json:"runtimeid"`
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
	FunctionID  int    `json:"functionid"`
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

	db, err := sql.Open("sqlite", filepath.Join(DBPath, "faas.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	handler.db = db

	if err := applyMigration(db); err != nil {
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

// GetRuntimeByName retrieves a runtime by name
func (h *DBHandler) GetRuntimeByName(name string) (*Runtime, error) {
	row := h.db.QueryRow("SELECT id, name, lang, repo, path, created_at, commitid FROM Runtimes WHERE name = ?", name)
	var runtime Runtime
	if err := row.Scan(&runtime.ID, &runtime.Name, &runtime.Lang, &runtime.Repo, &runtime.Path, &runtime.CreatedAt, &runtime.CommitID); err != nil {
		return nil, fmt.Errorf("failed to get runtime %s: %v", name, err)
	}
	return &runtime, nil
}

// ListRuntimes retrieves all runtimes
func (h *DBHandler) GetRuntimesList() ([]Runtime, error) {
	rows, err := h.db.Query("SELECT id, name, lang, repo, path, created_at, commitid FROM Runtimes")
	if err != nil {
		return nil, fmt.Errorf("failed to query runtimes: %v", err)
	}
	defer rows.Close()

	var runtimes []Runtime
	for rows.Next() {
		var r Runtime
		if err := rows.Scan(&r.ID, &r.Name, &r.Lang, &r.Repo, &r.Path, &r.CreatedAt, &r.CommitID); err != nil {
			return nil, fmt.Errorf("failed to scan runtime: %v", err)
		}
		runtimes = append(runtimes, r)
	}
	return runtimes, nil
}

// InsertFunction inserts a new function
func (h *DBHandler) InsertFunction(f Function) error {
	runtime, err := h.GetRuntimeByName(f.Runtime)
	if err != nil {
		return fmt.Errorf("failed to get runtime for function %s: %v", f.Name, err)
	}
	f.RuntimeID = runtime.ID

	_, err = h.db.Exec("INSERT INTO Functions (name, hash, description, runtime, guildid, runtimeid) VALUES (?, ?, ?, ?, ?, ?)",
		f.Name, f.Hash, f.Description, f.Runtime, f.GuildID, f.RuntimeID)
	if err != nil {
		return fmt.Errorf("failed to insert function %s: %v", f.Name, err)
	}
	return nil
}

// GetFunctionsByGuild retrieves functions for a given guild
func (h *DBHandler) GetFunctionsByGuild(guildID int) ([]Function, error) {
	rows, err := h.db.Query("SELECT id, name, hash, description, runtime, created_at, updated_at, guildid FROM Functions WHERE guildid = ?", guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to query functions for guild %d: %v", guildID, err)
	}
	defer rows.Close()

	var functions []Function
	for rows.Next() {
		var f Function
		if err := rows.Scan(&f.ID, &f.Name, &f.Hash, &f.Description, &f.Runtime, &f.CreatedAt, &f.UpdatedAt, &f.GuildID); err != nil {
			return nil, fmt.Errorf("failed to scan function for guild %d: %v", guildID, err)
		}
		functions = append(functions, f)
	}
	return functions, nil
}

// InsertCommand inserts a new command
func (h *DBHandler) InsertCommand(c Command) error {
	_, err := h.db.Exec("INSERT INTO Commands (command, description, guildid, functionid) VALUES (?, ?, ?, ?)",
		c.Command, c.Description, c.GuildID, c.FunctionID)
	if err != nil {
		return fmt.Errorf("failed to insert command %s: %v", c.Command, err)
	}
	return nil
}

// GetCommandsByGuild retrieves commands for a guild
func (h *DBHandler) GetCommandsByGuild(guildID int) ([]Command, error) {
	rows, err := h.db.Query("SELECT id, command, description, guildid, functionid FROM Commands WHERE guildid = ?", guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to query commands for guild %d: %v", guildID, err)
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.Command, &c.Description, &c.GuildID, &c.FunctionID); err != nil {
			return nil, fmt.Errorf("failed to scan command for guild %d: %v", guildID, err)
		}
		commands = append(commands, c)
	}
	return commands, nil
}

// InsertCommandArgument inserts a new command argument
func (h *DBHandler) InsertCommandArgument(arg CommandArgument) error {
	_, err := h.db.Exec("INSERT INTO CommandArguments (command_id, argument, description) VALUES (?, ?, ?)",
		arg.CommandID, arg.Name, arg.Description)
	if err != nil {
		return fmt.Errorf("failed to insert argument for command %d: %v", arg.CommandID, err)
	}
	return nil
}

// GetArgumentsByCommand retrieves arguments for a command
func (h *DBHandler) GetArgumentsByCommand(commandID int) ([]CommandArgument, error) {
	rows, err := h.db.Query("SELECT id, command_id, argument, description FROM CommandArguments WHERE command_id = ?", commandID)
	if err != nil {
		return nil, fmt.Errorf("failed to query arguments for command %d: %v", commandID, err)
	}
	defer rows.Close()

	var args []CommandArgument
	for rows.Next() {
		var arg CommandArgument
		if err := rows.Scan(&arg.ID, &arg.CommandID, &arg.Name, &arg.Description); err != nil {
			return nil, fmt.Errorf("failed to scan argument for command %d: %v", commandID, err)
		}
		args = append(args, arg)
	}
	return args, nil
}

// InsertRuntime inserts a new runtime
func (h *DBHandler) InsertRuntime(r Runtime) error {
	_, err := h.db.Exec("INSERT INTO Runtimes (name, lang, repo, path, commitid) VALUES (?, ?, ?, ?, ?)",
		r.Name, r.Lang, r.Repo, r.Path, r.CommitID)
	if err != nil {
		return fmt.Errorf("failed to insert runtime %s: %v", r.Name, err)
	}
	return nil
}
