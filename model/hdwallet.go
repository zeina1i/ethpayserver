package model

type HDWallet struct {
	ID   uint32 `json:"-" db:"id"`
	XPub string `json:"x_pub" db:"x_pub"`
}
