package model

type Merchant struct {
	ID       uint32 `json:"-" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string ``
}
