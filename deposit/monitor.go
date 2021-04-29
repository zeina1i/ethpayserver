package deposit

import (
	"github.com/zeina1i/ethpay/model"
	"github.com/zeina1i/ethpay/store"
	"github.com/zeina1i/ethpay/types"
)

type Monitor struct {
	store store.Store
}

func NewMonitor(store store.Store) *Monitor {
	return &Monitor{store: store}
}

func (m *Monitor) getNewBlockchainTxHandlerFunc() types.EventHandlerFunc {
	return func(event *types.Event) error {
		data := event.Data.(*types.NewBlockchainTxData)

		_, err := m.store.GetAddress(data.ToAddress)
		if err != nil {
			return nil
		}

		_, err = m.store.AddTx(&model.Tx{
			TxTime:      data.BlockTime,
			ReflectTime: data.BlockTime,
			FromAddress: data.FromAddress,
			ToAddress:   data.ToAddress,
			Asset:       data.Asset,
			Amount:      data.Amount,
			BlockNo:     data.BlockNo,
			TxHash:      data.TxHash,
			IsReflected: 1,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
