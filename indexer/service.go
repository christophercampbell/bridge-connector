package indexer

import (
	"github.com/christophercampbell/bridge-connector/config"
	"github.com/christophercampbell/bridge-connector/db"
)

type Service struct {
	config config.ChainConfig
	store  *db.Storage
}

func New(config config.ChainConfig, store *db.Storage) (*Service, error) {
	service := Service{
		config: config,
		store:  store,
	}
	// create client, do some sanity checks
	return &service, nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() {

}
