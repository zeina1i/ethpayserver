package model

type Merchant struct {
	ID       uint32 `json:"-" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string ``
}
