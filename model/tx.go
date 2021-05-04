package model

import "time"

type Tx struct {
	ID          uint32    `json:"-" db:"id"`
	TxTime      time.Time `json:"tx_time" db:"tx_time"`
	ReflectTime time.Time `json:"reflect_time" db:"reflect_time"`
	FromAddress string    `json:"from" db:"from_address"`
	ToAddress   string    `json:"to_address" db:"to_address"`
	Asset       string    `json:"asset" db:"asset"`
	Amount      float64   `json:"amount" db:"amount"`
	BlockNo     int64     `json:"block_no" db:"block_no"`
	TxHash      string    `json:"tx_hash" db:"tx_hash"`
	IsReflected uint      `json:"is_reflected" db:"is_reflected"`
}
