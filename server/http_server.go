package server

import (
	"context"
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/zeina1i/ethpay/config"
	"github.com/zeina1i/ethpay/hdwallet"
	"github.com/zeina1i/ethpay/httputil"
	"github.com/zeina1i/ethpay/model"
	"github.com/zeina1i/ethpay/passwords"
	"github.com/zeina1i/ethpay/store"
	"github.com/zeina1i/ethpay/types"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type HTTPServer struct {
	*httputil.Server
	templates *httputil.Templates

	store  store.Store
	pm     passwords.Passwords
	config config.Config
}

type ContextKey int

const (
	TokenContextKey ContextKey = iota
	MerchantContextKey
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewHTTPServer(store store.Store, pm passwords.Passwords) *HTTPServer {
	s := &HTTPServer{
		Server:    httputil.NewServer(httputil.NewRouter(), httputil.NewConfig(8082)),
		templates: httputil.NewTemplates("base"),

		store: store,
		pm:    pm,
	}

	s.Server.Router.ServeFiles(
		"/css/*filepath",
		rice.MustFindBox("../static/css").HTTPBox(),
	)
	s.Server.Router.ServeFiles(
		"/js/*filepath",
		rice.MustFindBox("../static/js").HTTPBox(),
	)
	s.Server.Router.ServeFiles(
		"/images/*filepath",
		rice.MustFindBox("../static/images").HTTPBox(),
	)

	return s
}

func (s *HTTPServer) InitRoutes() {
	s.Router.GET("/api/v1/ping", s.PingEndpoint())
	s.Router.POST("/api/v1/auth", s.AuthEndpoint())
	s.Router.POST("/api/v1/register", s.RegisterEndpoint())

	s.Router.POST("/api/v1/generate-address", s.GenerateAddress())
	s.Router.GET("/api/v1/transactions", s.isAuthorized(s.TransactionsEndpoint()))
}

func (s *HTTPServer) CreateToken(merchant *model.Merchant, r *http.Request) (*types.Token, error) {
	claims := jwt.MapClaims{}
	claims["email"] = merchant.Email
	createdAt := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.APISigningKey))
	if err != nil {
		return nil, err
	}

	signedToken, err := jwt.Parse(tokenString, s.jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	tkn := &types.Token{
		Signature: signedToken.Signature,
		Value:     tokenString,
		UserAgent: r.UserAgent(),
		CreatedAt: createdAt,
	}

	return tkn, nil
}

func (s *HTTPServer) jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}
	return []byte(s.config.APISigningKey), nil
}

func (s *HTTPServer) isAuthorized(endpoint httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Header.Get("Token") == "" {
			data, _ := json.Marshal(types.NewResponse(true, "No Token has provided", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)
			return
		}

		token, err := jwt.Parse(r.Header.Get("Token"), s.jwtKeyFunc)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Token is not valid", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			email := claims["email"].(string)

			merchant, err := s.store.GetMerchant(email)
			if err != nil {
				data, _ := json.Marshal(types.NewResponse(true, "Merchant not found", nil))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(data)
				return
			}

			ctx := context.WithValue(r.Context(), TokenContextKey, token)
			ctx = context.WithValue(ctx, MerchantContextKey, merchant)

			endpoint(w, r.WithContext(ctx), p)
		} else {
			data, _ := json.Marshal(types.NewResponse(true, "Token is not valid", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)
			return
		}
	}
}

func (s *HTTPServer) PingEndpoint() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
		return
	}
}

func (s *HTTPServer) AuthEndpoint() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req, err := types.NewAuthRequest(r.Body)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Bad Request", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		if req.Email == "" || req.Password == "" {
			data, _ := json.Marshal(types.NewResponse(true, "Bad Request", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		merchant, err := s.store.GetMerchant(req.Email)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Invalid credentials", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)
			return
		}

		err = s.pm.CheckPassword(merchant.Password, req.Password)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Invalid credentials", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)
			return
		}

		token, err := s.CreateToken(merchant, r)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Error creating token", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		data, _ := json.Marshal(types.AuthResponse{
			Token: token.Value,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
}

func (s *HTTPServer) RegisterEndpoint() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		req, err := types.NewRegisterRequest(r.Body)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Bad Request", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		email := strings.TrimSpace(req.Email)
		if !isEmailValid(email) {
			data, _ := json.Marshal(types.NewResponse(true, "Email is not valid", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		_, err = s.store.GetMerchant(req.Email)
		if err == nil {
			data, _ := json.Marshal(types.NewResponse(true, "Merchant exists", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}

		hash, err := s.pm.CreatePassword(req.Password)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Password creation failed", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
		}

		_, err = s.store.AddMerchant(&model.Merchant{
			Email:    req.Email,
			Password: hash,
		})
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Storing merchant failed", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		data, _ := json.Marshal(types.NewResponse(false, "Registration successful", nil))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
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

func (s *HTTPServer) TransactionsEndpoint() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		paginationItems := s.config.PaginationItems
		merchant := r.Context().Value(MerchantContextKey).(*model.Merchant)
		page := safeParseInt(r.FormValue("p"), 1)

		txs, err := s.store.GetTxs(merchant.ID, (page-1)*paginationItems, paginationItems)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Getting transaction list error", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		count, err := s.store.CountTxs(merchant.ID)
		if err != nil {
			data, _ := json.Marshal(types.NewResponse(true, "Counting transactions error", nil))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(data)
			return
		}

		var responseTxs []*types.ListTx
		for _, tx := range txs {
			responseTxs = append(responseTxs, &types.ListTx{
				ID:          tx.ID,
				TxTime:      tx.TxTime,
				ReflectTime: tx.ReflectTime,
				FromAddress: tx.FromAddress,
				ToAddress:   tx.ToAddress,
				Asset:       tx.Asset,
				Amount:      tx.Amount,
				BlockNo:     tx.BlockNo,
				TxHash:      tx.TxHash,
				IsReflected: tx.IsReflected,
			})
		}

		pagedResponse := types.TxsPagedResponse{
			Transactions: responseTxs,
			Pager: types.PagerResponse{
				Current:  page,
				MaxPages: count/paginationItems + 1,
				Total:    count,
			},
		}

		data, _ := json.Marshal(pagedResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
}
