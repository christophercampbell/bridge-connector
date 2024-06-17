package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/christophercampbell/bridge-connector/log"
	"github.com/christophercampbell/bridge-connector/types"
	"github.com/umbracle/ethgo"
)

// TODO: consider using context versions of DB queries/exec
type Storage struct {
	db *sql.DB
}

func NewStorage(file string) (*Storage, error) {
	var (
		db  *sql.DB
		err error
	)
	if err = os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		log.Errorf("invalid file path: %v", file)
		return nil, err
	}
	if db, err = ConnectSqlite3(file); err != nil {
		return nil, err
	}
	store := &Storage{
		db: db,
	}
	if err = store.init(); err != nil {
		store.Close()
		return nil, err
	}
	return store, nil
}

func (s *Storage) StoreEvents(chainId uint, events []types.BridgeEvent) error {
	if len(events) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, e := range events {
		_, err = tx.Exec(upsertEventStatement,
			chainId,
			e.BlockNumber,
			e.TransactionIndex,
			e.LogIndex,
			e.TransactionHash.String(),
			e.EventType,
			e.Removed)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = errors.Join(err, rbErr)
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ReadEvents(chainId uint, fromBlock, toBlock uint64) ([]types.BridgeEvent, error) {
	rows, err := s.db.Query(selectEventsStatement, chainId, fromBlock, toBlock)
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

func (s *Storage) UpdateLastProcessedBlock(chainId uint, block uint64) error {
	_, err := s.db.Exec(upsertLastBlockStatement, chainId, block)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetLastProcessedBlock(chainId uint) (uint64, error) {
	row := s.db.QueryRow(selectLastProcessedBlockStatement, chainId)
	var blockNumber uint64
	err := row.Scan(&blockNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil // TODO: instead of 0, it should be the contract's deployed block
	}
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (s *Storage) Close() {
	_ = s.db.Close()
}

// TODO: real migrations etc
func (s *Storage) init() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, ddl := range []string{eventTableDDL, lastProcessedBlockDDL} {
		if _, err = tx.Exec(ddl); err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = errors.Join(rbErr, err)
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
