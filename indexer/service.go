package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/christophercampbell/bridge-connector/config"
	"github.com/christophercampbell/bridge-connector/db"
	"github.com/christophercampbell/bridge-connector/log"
	"github.com/christophercampbell/bridge-connector/types"
	"github.com/pkg/errors"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

type Service struct {
	chainId    uint
	store      *db.Storage
	rpc        *jsonrpc.Client
	batchSize  uint
	once       sync.Once
	cancelFunc context.CancelFunc
	lastBlock  uint64
	rateLimit  time.Duration
}

func New(config config.ChainConfig, store *db.Storage) (*Service, error) {
	jsonrpcClient, err := jsonrpc.NewClient(config.RpcURL)
	if err != nil {
		return nil, err
	}
	chainId, err := jsonrpcClient.Eth().ChainID()
	if err != nil {
		return nil, err
	}
	if uint64(config.ChainId) != chainId.Uint64() {
		return nil,
			errors.Errorf("wrong chain id for %v, expected %d, got %d",
				config.RpcURL, config.ChainId, chainId.Uint64())
	}
	service := Service{
		chainId:    config.ChainId,
		store:      store,
		rpc:        jsonrpcClient,
		batchSize:  config.IndexerConfig.BlockBatchSize,
		cancelFunc: nil,
		once:       sync.Once{},
		rateLimit:  config.IndexerConfig.RateLimit.Duration,
	}

	var lastBlock uint64
	if lastBlock, err = store.GetLastProcessedBlock(config.ChainId); err != nil { // should this be last finalized block?
		return nil, err
	}
	if lastBlock == 0 {
		lastBlock = config.IndexerConfig.GenesisBlock
	}
	service.lastBlock = lastBlock

	return &service, nil
}

func (s *Service) Start(parentContext context.Context) error {
	s.once.Do(func() {
		ctx, cancel := context.WithCancel(parentContext)
		s.cancelFunc = cancel

		go s.processEvents(ctx)
	})
	return nil
}

func (s *Service) processEvents(ctx context.Context) {
	handleErrs := make(chan error)
	for {
		select {
		case err := <-handleErrs:
			log.Errorf("could not process events: %+v", err)
			<-time.After(1 * time.Second) // TODO: exponential backoff
		case <-ctx.Done():
			return
		default:
		}

		// basic rate limiter
		<-time.After(s.rateLimit)

		next := s.lastBlock + 1
		events, err := s.retrieveEvents(next, s.batchSize)
		if err != nil {
			handleErrs <- err
			continue
		}

		err = s.store.StoreEvents(s.chainId, events)
		if err != nil {
			handleErrs <- err
			continue
		}

		err = s.store.UpdateLastProcessedBlock(s.chainId, next+uint64(s.batchSize))
		if err != nil {
			handleErrs <- err
			continue
		}

		s.lastBlock = next + uint64(s.batchSize)
	}
}

func (s *Service) retrieveEvents(startBlock uint64, count uint) ([]types.BridgeEvent, error) {
	from := ethgo.BlockNumber(startBlock)
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
	for _, le := range logs {
		event := maybeFromLog(le)
		if event == nil {
			continue
		}
		events = append(events, *event)
	}
	return events, nil
}

func (s *Service) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
