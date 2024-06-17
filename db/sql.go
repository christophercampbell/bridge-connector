package db

const (
	enableWALSql string = "PRAGMA journal_mode=WAL;"

	// Event table likely needs a chainID? and the data blob, or maybe the event in the blob gets parsed and into a typed table
	eventTableDDL string = `
CREATE TABLE IF NOT EXISTS events (
    chain_id	 INTEGER,
  	block_number INTEGER,               
	tx_index 	 INTEGER,               
	log_index    INTEGER,               
	tx_hash  	 TEXT,
    event_type   INTEGER,
    removed      INTEGER,
    PRIMARY KEY (chain_id, block_number, tx_index)                  
);
CREATE INDEX IF NOT EXISTS events_chain_id_idx ON events(chain_id);
CREATE INDEX IF NOT EXISTS events_event_type_idx ON events(event_type);
`

	upsertEventStatement string = `
INSERT INTO events 
    (chain_id, block_number, tx_index, log_index, tx_hash, event_type, removed) 
VALUES 
    (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT DO UPDATE 
    SET log_index = EXCLUDED.log_index, tx_hash = EXCLUDED.tx_hash, event_type = EXCLUDED.event_type, removed = EXCLUDED.removed`

	selectEventsStatement string = `
SELECT block_number, tx_index, log_index, tx_hash, event_type, removed 
FROM events
	WHERE chain_id = ? 
	  AND block_number BETWEEN ? AND ?` // BETWEEN is inclusive on both ends

	lastProcessedBlockDDL string = `
CREATE TABLE IF NOT EXISTS last_processed_block (
    chain_id INTEGER PRIMARY KEY,
	block_number INTEGER NOT NULL
);
`
	upsertLastBlockStatement = `
INSERT INTO last_processed_block (chain_id, block_number) VALUES (?, ?) 
                                                          ON CONFLICT(chain_id)
                                                          DO UPDATE set block_number = EXCLUDED.block_number`

	selectLastProcessedBlockStatement = `SELECT block_number FROM last_processed_block WHERE chain_id = ?`
)
