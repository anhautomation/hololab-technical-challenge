package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, _ = db.Exec(`PRAGMA journal_mode=WAL;`)
	if err := initSchema(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS bots (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      job TEXT NOT NULL,
      bio TEXT NOT NULL,
      style TEXT NOT NULL,
      knowledge TEXT NOT NULL DEFAULT '',
      created_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      bot_id TEXT NOT NULL,
      role TEXT NOT NULL CHECK (role IN ('user','assistant')),
      content TEXT NOT NULL,
      created_at TEXT NOT NULL,
      FOREIGN KEY(bot_id) REFERENCES bots(id) ON DELETE CASCADE
    );

    CREATE INDEX IF NOT EXISTS idx_messages_bot_time
    ON messages (bot_id, created_at);
  `)
	return err
}
