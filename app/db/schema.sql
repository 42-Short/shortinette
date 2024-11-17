PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS % s (
  intra_login TEXT PRIMARY KEY NOT NULL UNIQUE,
  github_login TEXT NOT NULL UNIQUE,
  FOREIGN KEY (intra_login) REFERENCES participants(intra_login) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS % s (
  module_id INTEGER NOT NULL,
  intra_login TEXT NOT NULL,
  attempts INTEGER DEFAULT 0,
  score INTEGER DEFAULT 0,
  last_graded DATETIME,
  wait_time INTEGER DEFAULT 0,
  grading_ongoing BOOLEAN DEFAULT 0,
  PRIMARY KEY (module_id, intra_login),
  FOREIGN KEY (intra_login) REFERENCES participants(intra_login) ON DELETE CASCADE
);
