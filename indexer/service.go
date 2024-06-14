package indexer

import (
	"context"

	"github.com/christophercampbell/bridge-connector/config"
	"github.com/christophercampbell/bridge-connector/db"
	"github.com/christophercampbell/bridge-connector/types"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

type Service struct {
	config config.ChainConfig
	store  *db.Storage
	rpc    *jsonrpc.Client
}

func New(config config.ChainConfig, store *db.Storage) (*Service, error) {
	jsonrpcClient, err := jsonrpc.NewClient(config.RpcURL)
	if err != nil {
		return nil, err
	}

	service := Service{
		config: config,
		store:  store,
		rpc:    jsonrpcClient,
	}

	// create client, do some sanity checks
	return &service, nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) retrieveEvents(ctx context.Context, start, count uint64) ([]types.BridgeEvent, error) {
	from := ethgo.BlockNumber(start)
	to := from + ethgo.BlockNumber(count)
	filter := ethgo.LogFilter{
		// Address:   []ethgo.Address{ethgo.HexToAddress(s.config.)},
		BlockHash: nil,
		From:      &from,
		To:        &to,
	}
	logs, err := s.rpc.Eth().GetLogs(&filter)
	if err != nil {
		return nil, err
	}
	var events []types.BridgeEvent
	for _, log := range logs {
		// process log event
		event, err := s.parseEvent(log)
		if err != nil {
			// skip or fail?
			return nil, err
		}
		events = append(events, *event)
	}
	return events, nil
}

func (s *Service) Stop() {

}

func (s *Service) parseEvent(log *ethgo.Log) (*types.BridgeEvent, error) {
	return nil, nil
}
