package config

// Config represents the full configuration
type Config struct {
	LX ChainConfig
	LY ChainConfig
	DB DBConfig
}

// ChainConfig is a struct that defines contract and service settings on a particular chain
type ChainConfig struct {
	ChainId       uint          `mapstructure:"ChainId"`
	RpcURL        string        `mapstructure:"RpcURL"`
	BridgeAddress string        `mapstructure:"BridgeAddress"`
	IndexerConfig IndexerConfig `mapstructure:"Indexer"`
}

type IndexerConfig struct {
	Timeout        Duration `mapstructure:"Timeout"`
	RetryPeriod    Duration `mapstructure:"RetryPeriod"`
	BlockBatchSize uint     `mapstructure:"BlockBatchSize"`
	BridgeAddress  string   `mapstructure:"BridgeAddress"`
	GenesisBlock   uint64   `mapstructure:"GenesisBlock"`
}

type DBConfig struct {
	File string `mapstructure:"File"`
}
