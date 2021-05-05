package mysql

import "github.com/zeina1i/ethpay/model"

const (
	AddAddressStmt = `
INSERT INTO addresses
(address, hd_wallet_id, account_id, account_index, path, created_at)
VALUES (:address, :hd_wallet_id, :account_id, :account_index, :path, :created_at)`

	GetAddressStmt = `
SELECT * FROM addresses AS a
WHERE a.address=?;
`
)

func (s Store) AddAddress(address *model.Address) (*model.Address, error) {
	m := map[string]interface{}{
		"address":       address.Address,
		"hd_wallet_id":  address.HDWalletID,
		"account_id":    address.AccountId,
		"account_index": address.AccountIndex,
		"path":          address.Path,
		"created_at":    address.CreatedAt,
	}

	res, err := s.NamedExec(AddAddressStmt, m)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	address.ID = uint32(id)

	return address, nil
}

func (s *Store) GetAddress(addressHex string) (*model.Address, error) {
	var address model.Address

	err := s.Get(&address, GetAddressStmt, addressHex)
	if err != nil {
		return nil, err
	}

	return &address, nil
}
