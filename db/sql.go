package db

const (
	enableWALSql string = "PRAGMA journal_mode=WAL;"

	eventTableSql string = `CREATE TABLE IF NOT EXISTS events (
  	block_number INTEGER,               
	tx_index 	 INTEGER,               
	log_index    INTEGER,               
	tx_hash  	 TEXT,
    event_type   INTEGER,
    removed      INTEGER,
    PRIMARY KEY (block_number, tx_index)                  
  );`

	insertEventStatement string = `INSERT INTO events 
    (block_number, tx_index, log_index, tx_hash, event_type, removed) VALUES 
    (?, ?, ?, ?, ?, ?, ?)`

	selectEventsStatement string = `SELECT block_number, tx_index, log_index, tx_hash, event_type, removed FROM events
	WHERE block_number BETWEEN ? AND ?` // BETWEEN is inclusive on both ends
)
