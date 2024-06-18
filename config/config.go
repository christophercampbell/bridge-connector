package config

import (
	"github.com/christophercampbell/bridge-connector/log"
	"github.com/umbracle/ethgo"
)

// Config represents the full configuration
type Config struct {
	Contracts Contracts
	Chains    []ChainConfig
	DB        DBConfig
	Log       log.Config
}

// ChainConfig is a struct that defines contract and service settings on a particular chain
type ChainConfig struct {
	Name          string        `mapstructure:"Name"`
	Enabled       bool          `mapstructure:"Enabled"`
	ChainId       uint          `mapstructure:"ChainId"`
	RpcURL        string        `mapstructure:"RpcURL"`
	IndexerConfig IndexerConfig `mapstructure:"Indexer"`
}

type IndexerConfig struct {
	Timeout        Duration `mapstructure:"Timeout"`
	RateLimit      Duration `mapstructure:"RateLimit"`
	BlockBatchSize uint     `mapstructure:"BlockBatchSize"`
	GenesisBlock   uint64   `mapstructure:"GenesisBlock"`
}

type DBConfig struct {
	File string `mapstructure:"File"`
}

type Contracts struct {
	BridgeEthMainnetAddr string `mapstructure:"BridgeEthMainnetAddr"`
	GlobalExitRootAddr   string `mapstructure:"GlobalExitRootAddr"`
	RollupManagerAddr    string `mapstructure:"RollupManagerAddr"`
}

func (c *Contracts) Addresses() []ethgo.Address {
	return []ethgo.Address{
		ethgo.HexToAddress(c.BridgeEthMainnetAddr),
		ethgo.HexToAddress(c.GlobalExitRootAddr),
		ethgo.HexToAddress(c.RollupManagerAddr)}
}
