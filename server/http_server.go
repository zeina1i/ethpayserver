package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zeina1i/ethpay/hdwallet"
	"github.com/zeina1i/ethpay/httputil"
	"github.com/zeina1i/ethpay/model"
	"github.com/zeina1i/ethpay/store"
	"github.com/zeina1i/ethpay/types"
	"net/http"
	"time"
)

type HTTPServer struct {
	*httputil.Server

	store store.Store
}

func (s *HTTPServer) GenerateAddress() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		req, err := types.NewGenerateAddressRequest(r.Body)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Bad Request", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		addressHex, err := hdwallet.GenerateAddress(req.XPub, req.Id, req.Index)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Error in generating address", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		address, err := s.store.GetAddress(addressHex)
		if err == nil {
			data, _ := json.Marshal(types.GenerateAddressResponse{
				Address:      address.Address,
				AccountId:    address.AccountId,
				AccountIndex: address.AccountIndex,
				Path:         "",
				CreatedAt:    address.CreatedAt,
				IsNew:        false,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		hdWallet, err := s.store.GetHDWallet(req.XPub)
		if err != nil {
			hdWallet, err = s.store.AddHDWallet(&model.HDWallet{
				XPub: req.XPub,
			})
		}

		address, err = s.store.AddAddress(&model.Address{
			HDWalletID:   hdWallet.ID,
			Address:      addressHex,
			AccountId:    req.Id,
			AccountIndex: req.Index,
			Path:         "",
			CreatedAt:    time.Now(),
		})
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Error in storing address", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		data, _ := json.Marshal(types.GenerateAddressResponse{
			Address:      address.Address,
			AccountId:    address.AccountId,
			AccountIndex: address.AccountIndex,
			Path:         "",
			CreatedAt:    address.CreatedAt,
			IsNew:        true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
}

func (s *HTTPServer) InitRoutes() {
	s.Router.POST("/api/v1/generate-address", s.GenerateAddress())
}
