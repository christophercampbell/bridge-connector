package indexer

import (
	"context"
	"math"
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
	name       string
	chainId    uint
	contracts  []ethgo.Address
	store      *db.Storage
	rpc        *jsonrpc.Client
	batchSize  uint
	once       sync.Once
	cancelFunc context.CancelFunc
	lastBlock  uint64
	rateLimit  time.Duration
}

func New(config config.ChainConfig, contracts config.Contracts, store *db.Storage) (*Service, error) {
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
		name:       config.Name,
		chainId:    config.ChainId,
		contracts:  contracts.Addresses(),
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
	for {
		select { // non-blocking selection
		case <-ctx.Done():
			return
		default:
		}

		// basic rate limiter
		<-time.After(s.rateLimit)

		next := s.lastBlock + 1
		events, err := s.retrieveEvents(next, s.batchSize)
		if err != nil {
			log.Errorf("could not retrieve events: %+v", err)
			continue
		}

		err = s.store.StoreEvents(s.chainId, events)
		if err != nil {
			log.Errorf("could not store events: %+v", err)
			continue
		}

		// Don't advance beyond last block
		chainEnd, err := s.rpc.Eth().BlockNumber()
		if err != nil {
			log.Errorf("failed to get chain block length: %+v", err)
			continue
		}
		youngestBlock := uint64(math.Min(float64(chainEnd), float64(next+uint64(s.batchSize))))

		err = s.store.UpdateLastProcessedBlock(s.chainId, youngestBlock)
		if err != nil {
			log.Errorf("could not update last processed block: %+v", err)
			continue
		}

		s.lastBlock = youngestBlock
	}
}

// metrics should be kept for block rate & event rate & estimation of sync status

func (s *Service) retrieveEvents(startBlock uint64, count uint) ([]types.BridgeEvent, error) {
	from := ethgo.BlockNumber(startBlock)
	to := from + ethgo.BlockNumber(count)
	filter := ethgo.LogFilter{
		Address:   s.contracts,
		BlockHash: nil,
		From:      &from,
		To:        &to,
	}

	logs, err := s.rpc.Eth().GetLogs(&filter)
	if err != nil {

		// Might have to handle this error: {"code":-32602,"message":"query returned more than 10000 results"}
		// Maybe the count parameter is adaptive?

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
