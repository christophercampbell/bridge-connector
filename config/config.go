package config

import "github.com/christophercampbell/bridge-connector/log"

// Config represents the full configuration
type Config struct {
	LX  ChainConfig
	LY  ChainConfig
	DB  DBConfig
	Log log.Config
}

// ChainConfig is a struct that defines contract and service settings on a particular chain
type ChainConfig struct {
	ChainId       uint          `mapstructure:"ChainId"`
	RpcURL        string        `mapstructure:"RpcURL"`
	IndexerConfig IndexerConfig `mapstructure:"Indexer"`
}

type IndexerConfig struct {
	Timeout        Duration `mapstructure:"Timeout"`
	RetryPeriod    Duration `mapstructure:"RetryPeriod"`
	BlockBatchSize uint     `mapstructure:"BlockBatchSize"`
	GenesisBlock   uint64   `mapstructure:"GenesisBlock"`
}

type DBConfig struct {
	File string `mapstructure:"File"`
}
