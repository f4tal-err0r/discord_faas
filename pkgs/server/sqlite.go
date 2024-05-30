package server

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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
		defaultchan TEXT NOT NULL,
	);
	
	CREATE TABLE IF NOT EXISTS ApprovedRoles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		roleid TEXT NOT NULL,
		guildid INTEGER,
		FOREIGN KEY (guildid) REFERENCES GuildMetadata(guildid)
	);
	
	CREATE TABLE IF NOT EXISTS Commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		hash TEXT NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		guildid INTEGER,
		FOREIGN KEY (guildid) REFERENCES GuildMetadata(guildid)
	);
	`
	_, err = db.Exec(CreateTablesDb)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}
