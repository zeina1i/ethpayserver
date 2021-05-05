package deposit

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeina1i/ethpay/client"
	"github.com/zeina1i/ethpay/model"
	"github.com/zeina1i/ethpay/store"
	"github.com/zeina1i/ethpay/types"
	"math/big"
	"sync"
	"time"
)

type Monitor struct {
	store     store.Store
	ethClient client.EthClientInterface

	scMap map[string]types.Erc20Token
}

func NewMonitor(store store.Store) *Monitor {
	return &Monitor{store: store}
}

func (m *Monitor) processBlock(block ethTypes.Block) {
	wg := sync.WaitGroup{}
	for _, tx := range block.Transactions() {
		go func() {
			defer wg.Done()
			m.handleNewBlockchainTx(tx, block)
		}()
	}
	wg.Wait()
}

func (m *Monitor) handleNewBlockchainTx(transaction *ethTypes.Transaction, block ethTypes.Block) error {
	handler := m.getNewBlockchainTxHandlerFunc()

	sc, ok := m.scMap[transaction.To().String()]
	if ok {
		tries := 0
	GetReceipt:
		receipt, err := m.ethClient.GetTxReceipt(context.Background(), transaction.Hash())
		if err != nil {
			if tries > 3 {
				fmt.Println(err, "after 3 tries")
			}
			tries++
			time.Sleep(5 * time.Second)
			goto GetReceipt
		} else {
			for _, receiptLog := range receipt.Logs {
				_, err := m.store.GetAddress(receiptLog.Topics[2].Hex())
				if err != nil {
					continue
				}

				from := common.HexToAddress(receiptLog.Topics[1].Hex())

				val, err := client.ContractAbi.Unpack("Transfer", receiptLog.Data)
				if err != nil {
					fmt.Println(err)
				}

				err = handler(&types.Event{
					Type: types.NewBlockchainTx,
					Data: &types.NewBlockchainTxData{
						TxHash:      transaction.Hash().String(),
						ToAddress:   receiptLog.Topics[2].Hex(),
						FromAddress: from.String(),
						BlockTime:   time.Unix(int64(block.Time()), 0),
						BlockNo:     block.Number().Int64(),
						Asset:       sc.Symbol,
						Amount:      client.Wei2Float(val[0].(*big.Int), sc.Decimals),
					},
				})
				if err != nil {

				}
			}
		}
	} else {
		_, err := m.store.GetAddress(transaction.To().String())
		if err != nil {
			return nil
		}

		err = handler(&types.Event{
			Type: types.NewBlockchainTx,
			Data: &types.NewBlockchainTxData{
				TxHash:      transaction.Hash().String(),
				ToAddress:   transaction.To().String(),
				FromAddress: "",
				BlockTime:   time.Unix(int64(block.Time()), 0),
				BlockNo:     block.Number().Int64(),
				Asset:       "ETH",
				Amount:      client.Wei2Float(transaction.Value(), 18),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
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
