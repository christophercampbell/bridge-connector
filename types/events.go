package types

import (
	"encoding/json"

	"github.com/umbracle/ethgo"
)

type BridgeEvent struct {
	Removed          bool                   `json:"removed"`
	BlockNumber      uint64                 `json:"block_number"`
	TransactionIndex uint64                 `json:"transaction_index"`
	LogIndex         uint64                 `json:"log_index"`
	TransactionHash  ethgo.Hash             `json:"transaction_hash"`
	EventType        uint8                  `json:"event_type"`
	Data             map[string]interface{} `json:"event_data"`
}

func (be *BridgeEvent) JsonData() (string, error) {
	bytes, err := json.Marshal(be.Data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (be *BridgeEvent) SetData(jsonStr string) error {
	var data map[string]interface{} // TODO: convert []bytes -> hex?
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return err
	}
	be.Data = data
	return nil
}
