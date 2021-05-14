package store

import "github.com/zeina1i/ethpay/model"

type Store interface {
	InitializeDB() error

	GetAddress(address string) (*model.Address, error)
	AddAddress(address *model.Address) (*model.Address, error)

	AddTx(tx *model.Tx) (*model.Tx, error)
	GetTxs(merchantId uint32, offset int, limit int) ([]*model.Tx, error)

	GetHDWallet(xPub string) (*model.HDWallet, error)
	AddHDWallet(hdWallet *model.HDWallet) (*model.HDWallet, error)

	GetMerchant(email string) (*model.Merchant, error)
	AddMerchant(merchant *model.Merchant) (*model.Merchant, error)
}
