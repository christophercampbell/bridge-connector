package db

const (
	createEventTable string = `
  CREATE TABLE IF NOT EXISTS event (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  description TEXT
  );`

	createMessageTable string = `
  CREATE TABLE IF NOT EXISTS message (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  description TEXT
  );`
)
