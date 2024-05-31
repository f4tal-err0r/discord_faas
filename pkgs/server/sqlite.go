package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	_ "github.com/mattn/go-sqlite3"
)

type GuildMetaRow struct {
	guildid     uint64
	name        string
	owner       string
	defaultchan string
}

type ApprovedRolesRow struct {
	roleid  string
	guildid uint64
}

type CommandsTableRow struct {
	command       string
	hash          string
	last_modified *time.Time
	guildid       uint64
}

func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	return db, nil
}

func InitDB(db *sql.DB) error {
	CreateTablesDb := `
	CREATE TABLE IF NOT EXISTS GuildMetadata (
		guildid INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		textchan TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS ApprovedRoles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		roleid TEXT NOT NULL,
		guildid INTEGER,
		FOREIGN KEY(guildid) REFERENCES GuildMetadata(guildid)
	);
	
	CREATE TABLE IF NOT EXISTS Commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		hash TEXT NOT NULL,
		created TIMESTAMP NOT NULL,
		last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		guildid INTEGER,
		FOREIGN KEY(guildid) REFERENCES GuildMetadata(guildid)
	);
	`
	_, err := db.Exec(CreateTablesDb)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	return nil
}

func RowWriter(db *sql.DB, v interface{}) error {
	switch v.(type) {
	case GuildMetaRow:
		row := v.(GuildMetaRow)
		insertSQL := `INSERT INTO GuildMetadata (guildid, name, owner, defaultchan) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, row.guildid, row.name, row.owner, row.defaultchan)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	case CommandsTableRow:
		row := v.(CommandsTableRow)
		time := time.Now()
		insertSQL := `INSERT INTO ApprovedRoles (command, hash, last_modified, guildid) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, row.command, row.hash, time, row.guildid)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	case ApprovedRolesRow:
		row := v.(ApprovedRolesRow)
		insertSQL := `INSERT INTO ApprovedRoles (roleid, guildid) VALUES (?, ?)`
		_, err := db.Exec(insertSQL, row.roleid, row.guildid)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	}
	return fmt.Errorf("Error: Invalid type provided to Row Writer")
}

func InitGuildData(db *sql.DB) {
	var newGuilds []string

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("ERR: Unable to fetch config: %v", err)
	}

	session := GetSession(cfg)

	botGuilds, err := session.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("Error getting guilds: %v", err)
	}

	for _, guild := range botGuilds {
		var guildRow GuildMetaRow
		err := db.QueryRow("SELECT guildid FROM GuildMetadata WHERE guildid = ?", guild.ID).Scan(&guildRow.guildid)
		if err != nil {
			if err == sql.ErrNoRows {
				newGuilds = append(newGuilds, guild.ID)
			}
		}
	}

	for _, gid := range newGuilds {
		if err != nil {
			log.Fatalf("Error getting default channel: %v", err)
		}

		guild := GetGuildInfo(session, gid)

		defaultChan, err := GetDefaultChannel(session, guild.ID)
		if err != nil {
			log.Fatalf("Error getting default channel: %v", err)
		}

		guildid := strToUint(guild.ID)
		guildRow := GuildMetaRow{
			guildid:     guildid,
			name:        guild.Name,
			owner:       guild.OwnerID,
			defaultchan: defaultChan.ID,
		}
		_ = RowWriter(db, guildRow)
	}
}

func UpdateCmdsDb(db *sql.DB, cmd *CommandsTableRow) {
	sql := `UPDATE Commands SET last_modified = ?, hash = ? WHERE command = ? AND guildid = ?`
	time := time.Now()
	_, err := db.Exec(sql, time, cmd.hash, cmd.command, cmd.guildid)
	if err != nil {
		log.Fatalf("Error editing command: %v", err)
	}
}

func GetCmdsDb(db *sql.DB, cmd string, guildid uint64) *CommandsTableRow {
	sqlselect := `SELECT command, hash, last_modified, guildid FROM Commands WHERE command = ? AND guildid = ?`
	var row CommandsTableRow
	row.guildid = guildid
	err := db.QueryRow(sqlselect, cmd, guildid).Scan(&row.command, &row.hash, &row.last_modified, &row.guildid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalf("Error fetching command: %v", err)
	}
	return &row
}

// Remove guild row from db if bot is no longer in guild
func RemoveGuildRow(db *sql.DB, guildid string) {
	sql := `DELETE FROM GuildMetadata WHERE guildid = ?`
	_, err := db.Exec(sql, guildid)
	if err != nil {
		log.Fatalf("Error deleting guild: %v", err)
	}
}