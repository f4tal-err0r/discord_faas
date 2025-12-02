package db

import (
	"database/sql"
	"fmt"
)

func applyMigration(db *sql.DB) error {
	CreateTablesDb := `
		CREATE TABLE IF NOT EXISTS GuildMetadata (
			guildid INTEGER PRIMARY KEY,
			name TEXT NOT NULL CHECK (length(name) > 0),
			owner TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS Runtimes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL CHECK (length(name) > 0) UNIQUE,
			lang TEXT NOT NULL CHECK (length(lang) > 0),
			repo TEXT NOT NULL CHECK (length(repo) > 0),
			path TEXT NOT NULL CHECK (length(path) > 0),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			commitid TEXT NOT NULL CHECK (length(commitid) > 0)
		);

		CREATE TABLE IF NOT EXISTS Functions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL CHECK (length(name) > 0),
			description TEXT NOT NULL,
			runtime TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			guildid INTEGER NOT NULL,
			runtimeid INTEGER NOT NULL,
			FOREIGN KEY (runtimeid) REFERENCES Runtimes(id),
			FOREIGN KEY (guildid) REFERENCES GuildMetadata(guildid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS ApprovedRoles (
			roleid TEXT NOT NULL,
			guildid INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (roleid, guildid),
			FOREIGN KEY (guildid) REFERENCES GuildMetadata(guildid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS Commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			command TEXT NOT NULL CHECK (length(command) > 0),
			description TEXT NOT NULL,
			guildid INTEGER NOT NULL,
			functionid INTEGER NOT NULL,
			FOREIGN KEY (functionid) REFERENCES Functions(id) ON DELETE CASCADE,
			FOREIGN KEY (guildid) REFERENCES GuildMetadata(guildid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS CommandArguments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			command_id INTEGER NOT NULL,
			argument TEXT NOT NULL,
			description TEXT NOT NULL,
			FOREIGN KEY (command_id) REFERENCES Commands(id) ON DELETE CASCADE
		);
	`
	_, err := db.Exec(CreateTablesDb)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}
