package types

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(error bool, message string, data interface{}) *Response {
	return &Response{
		Error:   error,
		Message: message,
		Data:    data}
}

type GenerateAddressRequest struct {
	XPub  string `json:"x_pub"`
	Id    uint32 `json:"id"`
	Index uint32 `json:"index"`
}

func NewGenerateAddressRequest(r io.Reader) (req GenerateAddressRequest, err error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &req)
	return
}

type GenerateAddressResponse struct {
	Address      string    `json:"address"`
	AccountId    uint32    `json:"account_id"`
	AccountIndex uint32    `json:"account_index"`
	Path         string    `json:"path"`
	CreatedAt    time.Time `json:"created_at"`
	IsNew        bool      `json:"is_new"`
}

type RegisterRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewRegisterRequest(r io.Reader) (req RegisterRequest, err error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &req)
	return
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthRequest(r io.Reader) (req AuthRequest, err error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &req)
	return
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ListTx struct {
	ID          uint32    `json:"id"`
	TxTime      time.Time `json:"tx_time"`
	ReflectTime time.Time `json:"reflect_time"`
	FromAddress string    `json:"from"`
	ToAddress   string    `json:"to_address"`
	Asset       string    `json:"asset"`
	Amount      float64   `json:"amount"`
	BlockNo     int64     `json:"block_no"`
	TxHash      string    `json:"tx_hash"`
	IsReflected uint      `json:"is_reflected"`
}

type PagerResponse struct {
	Current  int `json:"current_page"`
	MaxPages int `json:"max_pages"`
	Total    int `json:"total"`
}

type TxsPagedResponse struct {
	Transactions []*ListTx `json:"transactions"`
	Pager        PagerResponse
}

type GenerateHDWalletRequest struct {
	XPub string `json:"x_pub"`
}

func NewGenerateHDWalletRequest(r io.Reader) (req GenerateHDWalletRequest, err error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &req)
	return
}

type ListHDWallet struct {
	ID   uint32 `json:"id"`
	XPub string `json:"x_pub"`
}
