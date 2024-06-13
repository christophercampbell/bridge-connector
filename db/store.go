package db

import (
	"database/sql"

	"github.com/christophercampbell/bridge-connector/types"
	"github.com/umbracle/ethgo"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(file string) (*Storage, error) {
	db, err := ConnectSqlite3(file)
	if err != nil {
		return nil, err
	}
	store := &Storage{
		db: db,
	}
	err = store.init()
	if err != nil {
		store.Close()
		return nil, err
	}
	return store, nil
}

func (s *Storage) InsertEvent(be types.BridgeEvent) error {
	_, err := s.db.Exec(insertEventStatement,
		be.BlockNumber, be.TransactionIndex, be.LogIndex, be.TransactionHash.String(), be.EventType, be.Removed)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadEvents(fromBlock, toBlock uint64) ([]types.BridgeEvent, error) {
	rows, err := s.db.Query(selectEventsStatement, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	var events []types.BridgeEvent
	for rows.Next() {
		var bn, tx, lx uint64
		var hash string
		var et uint8
		var r bool
		if err = rows.Scan(&bn, &tx, &lx, &hash, &et, &r); err != nil {
			return nil, err
		}
		e := types.BridgeEvent{
			Removed:          r,
			BlockNumber:      bn,
			TransactionIndex: tx,
			LogIndex:         lx,
			TransactionHash:  ethgo.HexToHash(hash),
			EventType:        et,
			Data:             nil, // handle data later
		}
		events = append(events, e)
	}
	return events, nil
}

func (s *Storage) Close() {
	_ = s.db.Close()
}

func (s *Storage) init() error {
	if _, err := s.db.Exec(eventTableSql); err != nil {
		return err
	}
	return nil
}
