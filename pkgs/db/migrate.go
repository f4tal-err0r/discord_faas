package db

import (
	"fmt"
)

func applyMigration(h *DBHandler) error {
	CreateTablesDb := `
		CREATE TABLE IF NOT EXISTS GuildMetadata (
			guildid INTEGER PRIMARY KEY,
			source TEXT NOT NULL,
			name TEXT NOT NULL CHECK (length(name) > 0),
			owner TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS Functions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL CHECK (length(name) > 0),
			description TEXT NOT NULL,
			runtime TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			guildid INTEGER NOT NULL,
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
	_, err := h.db.Exec(CreateTablesDb)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}
