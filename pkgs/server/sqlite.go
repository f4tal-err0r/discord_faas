package server

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/cache"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	_ "github.com/mattn/go-sqlite3"
)

type GuildMetaRow struct {
	Guildid  int64
	Source   string
	Name     string
	Owner    string
	Textchan string
}

type ApprovedRolesRow struct {
	Roleid  string
	Guildid int64
}

type CommandsTableRow struct {
	Command       string
	Hash          string
	Last_modified *time.Time
	Guildid       int64
	Desc          string
	Runtime       string
}

var CmdsCache = cache.New()
var RolesCache = cache.New()
var GuildMetaCache = cache.New()

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
		source TEXT NOT NULL,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		textchan TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS ApprovedRoles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		roleid TEXT NOT NULL,
		guildid INTEGER NOT NULL,
		FOREIGN KEY(guildid) REFERENCES GuildMetadata(guildid)
	);
	
	CREATE TABLE IF NOT EXISTS Commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		hash TEXT NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_modified TIMESTAMP NOT NULL,
		guildid INTEGER NOT NULL,
		desc TEXT,
		runtime TEXT NOT NULL,
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
	switch T := v.(type) {
	case GuildMetaRow:
		row := T
		insertSQL := `INSERT INTO GuildMetadata (guildid, source, name, owner, textchan) VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, row.Guildid, row.Source, row.Name, row.Owner, row.Textchan)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	case CommandsTableRow:
		row := T
		time := time.Now()
		insertSQL := `INSERT INTO Commands (command, hash, last_modified, guildid, desc, runtime) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(insertSQL, row.Command, row.Hash, time, row.Guildid, row.Desc, row.Runtime)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	case ApprovedRolesRow:
		row := T
		insertSQL := `INSERT INTO ApprovedRoles (roleid, guildid) VALUES (?, ?)`
		_, err := db.Exec(insertSQL, row.Roleid, row.Guildid)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %v", err)
		}
		return nil
	}
	return fmt.Errorf("Error: Invalid type provided to Row Writer")
}

func InitGuildData(session *discordgo.Session, db *sql.DB) error {
	var newGuilds []string

	botGuilds, err := session.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("Error getting guilds: %v", err)
	}

	for _, guild := range botGuilds {
		var guildRow GuildMetaRow
		err := db.QueryRow("SELECT guildid FROM GuildMetadata WHERE guildid = ?", guild.ID).Scan(&guildRow.Guildid)
		if err != nil {
			if err == sql.ErrNoRows {
				newGuilds = append(newGuilds, guild.ID)
			}
		}
	}

	for _, gid := range newGuilds {
		if err != nil {
			return fmt.Errorf("Error getting default channel: %v", err)
		}

		guild, err := GetGuildInfo(session, gid)
		if err != nil {
			return fmt.Errorf("Error getting guild: %v", err)
		}

		defaultChan, err := GetDefaultChannel(session, guild.ID)
		if err != nil {
			return fmt.Errorf("Error getting default channel: %v", err)
		}

		guildid, err := strconv.ParseInt(guild.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Error converting guildid: %v", err)
		}
		guildRow := GuildMetaRow{
			Guildid:  guildid,
			Name:     guild.Name,
			Owner:    guild.OwnerID,
			Textchan: defaultChan.ID,
		}
		_ = RowWriter(db, guildRow)
	}

	return nil
}

func UpdateCmdsDb(db *sql.DB, cmd *CommandsTableRow) {
	sql := `UPDATE Commands SET last_modified = ?, hash = ? WHERE command = ? AND guildid = ?`
	time := time.Now()
	_, err := db.Exec(sql, time, cmd.Hash, cmd.Command, cmd.Guildid)
	if err != nil {
		log.Fatalf("Error editing command: %v", err)
	}
}

func GetCmdsDb(db *sql.DB, guildid int64) ([]*CommandsTableRow, error) {
	// Check cache by guildid
	if cmds, ok := CmdsCache.Get(fmt.Sprintf("%d", guildid)); ok {
		return cmds.([]*CommandsTableRow), nil
	}

	sqlselect := `SELECT command, hash, last_modified, guildid, desc, runtime FROM Commands WHERE guildid = ?`
	rows, err := db.Query(sqlselect, guildid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []*CommandsTableRow
	for rows.Next() {
		cmd := &CommandsTableRow{}
		if err := rows.Scan(&cmd.Command, &cmd.Hash, &cmd.Last_modified, &cmd.Guildid, &cmd.Desc, &cmd.Runtime); err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}

	CmdsCache.Set(fmt.Sprintf("%d", guildid), cmds, 45*time.Minute)

	return cmds, nil
}

// Remove rows from all tables by guildid
func DeleteGuildData(db *sql.DB, guildid string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM GuildMetadata WHERE guildid = ?", guildid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM Commands WHERE guildid = ?", guildid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM ApprovedRoles WHERE guildid = ?", guildid)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get row from Roles table by guildid
func GetRolebyGuildid(db *sql.DB, guildid int64) ([]string, error) {
	// Check Roles slice by guildid
	if val, ok := RolesCache.Get(fmt.Sprintf("%d", guildid)); ok {
		return val.([]string), nil
	}

	sqlselect := `SELECT roleid FROM ApprovedRoles WHERE guildid = ?`
	rows, err := db.Query(sqlselect, guildid)
	if err != nil {
		return nil, fmt.Errorf("error fetching roles: %v", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("error fetching roles: %v", err)
		}
		roles = append(roles, role)
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles found for guildid: %d", guildid)
	}

	RolesCache.Set(fmt.Sprintf("%d", guildid), roles, 90*time.Minute)
	return roles, nil
}

// Add row to Roles table
func AddRole(db *sql.DB, guildid int64, roleid string) error {
	sql := `INSERT INTO ApprovedRoles (guildid, roleid) VALUES (?, ?)`
	_, err := db.Exec(sql, guildid, roleid)
	if err != nil {
		return fmt.Errorf("Error adding role: %v", err)
	}

	log.Printf("Added role %s to guild %d", roleid, guildid)
	return nil
}

// Lookup guild from guildid
func LookupGuild(db *sql.DB, guildid int64) (*GuildMetaRow, error) {
	// Check GuildMetaCache by guildid
	if val, ok := GuildMetaCache.Get(fmt.Sprintf("%d", guildid)); ok {
		return val.(*GuildMetaRow), nil
	}
	sqlselect := `SELECT guildid, source, name, owner, textchan FROM GuildMetadata WHERE guildid = ?`
	row := db.QueryRow(sqlselect, guildid)
	guild := &GuildMetaRow{}
	if err := row.Scan(&guild.Guildid, &guild.Source, &guild.Name, &guild.Owner, &guild.Textchan); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	GuildMetaCache.Set(fmt.Sprintf("%d", guildid), guild, 45*time.Minute)
	return guild, nil
}
