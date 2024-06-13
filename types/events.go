package types

import "github.com/umbracle/ethgo"

type BridgeEvent struct {
	Removed          bool                   `json:"removed"`
	BlockNumber      uint64                 `json:"block_number"`
	TransactionIndex uint64                 `json:"transaction_index"`
	LogIndex         uint64                 `json:"log_index"`
	TransactionHash  ethgo.Hash             `json:"transaction_hash"`
	EventType        uint8                  `json:"event_type"`
	Data             map[string]interface{} `json:"event_data"`
}
