package model

import "time"

type Address struct {
	ID           uint32    `json:"-" db:"id"`
	HDWalletID   uint32    `json:"hd_wallet_id" db:"hd_wallet_id"`
	HDWallet     *HDWallet `json:"hd_wallet"`
	Address      string    `json:"address" db:"address"`
	AccountId    uint32    `json:"account_id" db:"account_id"`
	AccountIndex uint32    `json:"account_index" db:"account_index"`
	Path         string    `json:"path" db:"path"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
