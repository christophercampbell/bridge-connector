package db

import "database/sql"

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

// Add methods for INSERT/UPDATE/SELECT etc...

func (s *Storage) Close() {
	_ = s.db.Close()
}

func (s *Storage) init() error {
	if _, err := s.db.Exec(createEventTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(createMessageTable); err != nil {
		return err
	}
	return nil
}
