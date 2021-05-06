package mysql

import "github.com/zeina1i/ethpay/model"

var (
	GetMerchantStmt = `
SELECT * FROM merchants AS m
WHERE m.email=?;
`

	AddMerchantStmt = `
INSERT INTO merchants
(email, password)
VALUES (:email, :password)
`
)

func (s *Store) GetMerchant(email string) (*model.Merchant, error) {
	var merchant model.Merchant

	err := s.Get(&merchant, GetMerchantStmt, email)
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (s *Store) AddMerchant(merchant *model.Merchant) (*model.Merchant, error) {
	m := map[string]interface{}{
		"email":    merchant.Email,
		"password": merchant.Password,
	}

	res, err := s.NamedExec(AddMerchantStmt, m)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	merchant.ID = uint32(id)

	return merchant, nil

}
