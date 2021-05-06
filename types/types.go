package types

import "time"

type Erc20Token struct {
	Symbol   string
	Decimals uint
	Address  string
}

type Token struct {
	Signature string
	Value     string
	UserAgent string
	CreatedAt time.Time
	ExpiresAt time.Time
}
