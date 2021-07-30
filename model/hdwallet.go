package model

type HDWallet struct {
	ID         uint32    `json:"-" db:"id"`
	XPub       string    `json:"x_pub" db:"x_pub"`
	MerchantID uint32    `json:"merchant_id" db:"merchant_id"`
	Merchant   *Merchant `json:"merchant"`
}
