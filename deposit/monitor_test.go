package deposit

import (
	"github.com/golang/mock/gomock"
	"github.com/zeina1i/ethpay/model"
	"github.com/zeina1i/ethpay/store/mock_store"
	"github.com/zeina1i/ethpay/types"
	"testing"
	"time"
)

func TestMonitor_getNewBlockchainTxHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	eventData := &types.NewBlockchainTxData{
		TxHash:      "0x2e4c19899ecdd019f19513b2d27ec8ad95a7609c9755f95775da0de031e43dbc",
		ToAddress:   "0xca3c9c0cf83b2eed0ed2af380e6e75a791b2d61b",
		FromAddress: "0x2fc0c3fa9345d60b968b759e856460a1229f99f1",
		BlockTime:   time.Now(),
		BlockNo:     8499316,
		Asset:       "ETH",
		Amount:      0.011816097216939852,
	}
	event := &types.Event{
		Type: types.NewBlockchainTx,
		Data: eventData,
	}

	store := mock_store.NewMockStore(ctrl)
	store.EXPECT().GetAddress("0xca3c9c0cf83b2eed0ed2af380e6e75a791b2d61b").Return(&model.Address{}, nil)
	store.EXPECT().AddTx(&model.Tx{
		TxTime:      eventData.BlockTime,
		ReflectTime: eventData.BlockTime,
		FromAddress: eventData.FromAddress,
		ToAddress:   eventData.ToAddress,
		Asset:       eventData.Asset,
		Amount:      eventData.Amount,
		BlockNo:     eventData.BlockNo,
		TxHash:      eventData.TxHash,
		IsReflected: 1,
	})
	monitor := NewMonitor(store)

	handler := monitor.getNewBlockchainTxHandlerFunc()
	if err := handler(event); err != nil {
		t.Fatalf("%v", err)
	}
}
