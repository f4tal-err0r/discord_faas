package funcstore

import "database/sql"

type Funcstore interface {
	Get(db *sql.DB, cmd string, gid string) (string, error)
	Create(db *sql.DB, cmd string, hash string, gid string) error
	Exec(cmd string, gid string) error
	Delete(db *sql.DB, cmd string, gid string) error
}
