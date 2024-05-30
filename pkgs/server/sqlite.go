package server

import (
	"database/sql"
	"fmt"
	"time"

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
	command string
	hash    string
	created *time.Time
	guildid uint64
}

func InitDB(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	CreateTablesDb := `
	CREATE TABLE IF NOT EXISTS GuildMetadata (
		guildid INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		defaultchan TEXT NOT NULL
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
	_, err = db.Exec(CreateTablesDb)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
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
		insertSQL := `INSERT INTO ApprovedRoles (command, hash, created, guildid) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, row.command, row.hash, row.created, row.guildid)
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
