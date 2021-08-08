package mysql

import "github.com/zeina1i/ethpay/model"

const (
	AddHDWalletStmt = `
INSERT INTO hd_wallets(x_pub, merchant_id)
VALUES (:x_pub, :merchant_id)
`
	GetHDWalletStmt = `
SELECT * FROM hd_wallets AS w
WHERE w.x_pub =?
`
	GetHDWalletsByMerchantIDStmt = `
SELECT * FROM hd_wallets AS w
WHERE w.merchant_id = :merchant_id
`
)

func (s *Store) AddHDWallet(wallet *model.HDWallet) (*model.HDWallet, error) {
	m := map[string]interface{}{
		"x_pub":       wallet.XPub,
		"merchant_id": wallet.MerchantID,
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

func (s *Store) GetHDWalletsByMerchantID(merchantID uint32) ([]*model.HDWallet, error) {
	m := map[string]interface{}{
		"merchant_id": merchantID,
	}

	var hdWallets []*model.HDWallet
	rows, err := s.NamedQuery(GetHDWalletsByMerchantIDStmt, m)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var hdWallet model.HDWallet
		err = rows.StructScan(&hdWallet)
		if err != nil {
			return nil, err
		}
		hdWallets = append(hdWallets, &hdWallet)
	}

	return hdWallets, err
}
