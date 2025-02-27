PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS participant (
  intra_login TEXT PRIMARY KEY NOT NULL UNIQUE,
  github_login TEXT NOT NULL UNIQUE,
  current_module_id INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS  module(
  id INTEGER NOT NULL,
  intra_login TEXT NOT NULL,
  attempts INTEGER DEFAULT 0,
  score INTEGER DEFAULT 0,
  last_graded DATETIME,
  wait_time INTEGER DEFAULT 0,
  PRIMARY KEY (id, intra_login),
  FOREIGN KEY (intra_login) REFERENCES participant(intra_login) ON DELETE CASCADE
);
