package mysql

import (
	"github.com/zeina1i/ethpay/model"
)

const AddTxStmt = `
INSERT INTO txs
(tx_time, reflect_time, from_address, to_address, asset, amount, block_no, tx_hash, is_reflected, tx_status)
VALUES (:tx_time, :reflect_time, :from_address, :to_address, :asset, :amount, :block_no, :tx_hash, :is_reflected, :tx_status)`

const GetTxsStmt = `
SELECT * FROM txs WHERE merchant_id=:merchant_id LIMIT :limit OFFSET :offset`

const CountTxsStmt = `
SELECT count(*) FROM txs WHERE merchant_id=:merchant_id`

func (s *Store) AddTx(tx *model.Tx) (*model.Tx, error) {
	m := map[string]interface{}{
		"tx_time":      tx.TxTime.Format("2006-01-02 15:04:05"),
		"reflect_time": tx.ReflectTime.Format("2006-01-02 15:04:05"),
		"from_address": tx.FromAddress,
		"to_address":   tx.ToAddress,
		"asset":        tx.Asset,
		"amount":       tx.Amount,
		"block_no":     tx.BlockNo,
		"tx_hash":      tx.TxHash,
		"is_reflected": tx.IsReflected,
	}

	res, err := s.NamedExec(AddTxStmt, m)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	tx.ID = uint32(id)

	return tx, nil
}

func (s *Store) GetTxs(merchantId uint32, offset int, limit int) ([]*model.Tx, error) {
	m := map[string]interface{}{
		"merchant_id": merchantId,
		"offset":      offset,
		"limit":       limit,
	}

	var txs []*model.Tx
	rows, err := s.NamedQuery(GetTxsStmt, m)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tx model.Tx
		err = rows.StructScan(&tx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, &tx)
	}

	return txs, err
}

func (s *Store) CountTxs(merchantId uint32) (int, error) {
	m := map[string]interface{}{
		"merchant_id": merchantId,
	}

	rows, err := s.NamedQuery(CountTxsStmt, m)
	if err != nil {
		return -1, err
	}

	var count int
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, err
}
