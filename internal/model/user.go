package model

const (
	StatusCreated string = "created"
	StatusDeleted string = "deleted"
)

type User struct {
	UserID      uint64
	ID          uint64  `json:"-"`
	Name        *string `json:"name"`  `json:"name"`
	PhoneNumber *string `json:"phone_number"`  `json:"phone_number"`
	Email       *string `json:"email"`  `json:"naemailme"`
	Raiting     float64 `json:"raiting"`
	Status      string
}
