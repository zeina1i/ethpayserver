package store

import "github.com/zeina1i/ethpay/model"

type Store interface {
	GetAddress(address string) (*model.Address, error)

	AddTx(tx *model.Tx) (*model.Tx, error)
}
