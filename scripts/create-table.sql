CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_users_upd AFTER UPDATE ON users
BEGIN
    UPDATE users
    SET update_at = DATETIME('now', 'localtime') 
    WHERE rowid == NEW.rowid;
END;