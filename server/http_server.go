package server

import (
	"encoding/json"
	"fmt"
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

	store  store.Store
	pm     passwords.Passwords
	config config.Config
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewHTTPServer(store store.Store) *HTTPServer {
	return &HTTPServer{
		Server: httputil.NewServer(httputil.NewRouter(), httputil.NewConfig(8082)),

		store: store,
	}
}

func (s *HTTPServer) InitRoutes() {
	s.Router.GET("/api/v1/ping", s.PingEndpoint())
	s.Router.POST("/api/v1/auth", s.PingEndpoint())
	s.Router.POST("/api/v1/register", s.AuthEndpoint())

	s.Router.POST("/api/v1/generate-address", s.GenerateAddress())
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

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
