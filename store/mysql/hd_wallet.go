package mysql

import "github.com/zeina1i/ethpay/model"

const (
	AddHDWalletStmt = `
INSERT INTO hd_wallets(x_pub)
VALUES (:x_pub)
`
	GetHDWalletStmt = `
SELECT * FROM hd_wallets AS w
WHERE w.x_pub =?
`
)

func (s *Store) AddHDWallet(wallet *model.HDWallet) (*model.HDWallet, error) {
	m := map[string]interface{}{
		"x_pub": wallet.XPub,
	}

	res, err := s.NamedExec(AddHDWalletStmt, m)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	wallet.ID = uint32(id)

	return wallet, nil
}

func (s *Store) GetHDWallet(xPub string) (*model.HDWallet, error) {
	var wallet model.HDWallet

	err := s.Get(&wallet, GetHDWalletStmt, xPub)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
