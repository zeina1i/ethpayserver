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
