package mysql

import "github.com/zeina1i/ethpay/model"

const AddTxStmt = `
INSERT INTO txs
(tx_time, reflect_time, from_address, to_address, asset, amount, block_no, tx_hash, is_reflected, tx_status)
VALUES (:tx_time, :reflect_time, :from_address, :to_address, :asset, :amount, :block_no, :tx_hash, :is_reflected, :tx_status)`

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
