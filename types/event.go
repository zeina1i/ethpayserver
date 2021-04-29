package types

import "time"

type EventType string

const NewBlockchainTx = "new_blockchain_tx"

type Event struct {
	Type EventType
	Data interface{}
}

type EventHandlerFunc func(event *Event) error

type NewBlockchainTxData struct {
	TxHash      string
	ToAddress   string
	FromAddress string
	BlockTime   time.Time
	BlockNo     int64
	Asset       string
	Amount      float64
}
