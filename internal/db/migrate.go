package db

import "database/sql"

func Init(db *sql.DB) error {

	schema := `
CREATE TABLE IF NOT EXISTS machines (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	ip TEXT,
	domain TEXT,
	role TEXT
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	password TEXT,
	hash TEXT
);
`

	_, err := db.Exec(schema)
	return err
}
