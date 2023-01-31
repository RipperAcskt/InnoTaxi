package model

type User struct {
	ID          uint64  `json:"-"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"naemailme"`
	Raiting     float64 `json:"raiting"`
}
