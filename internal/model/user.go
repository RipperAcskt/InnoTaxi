package model

const (
	StatusCreated string = "created"
	StatusDeleted string = "deleted"
)

type User struct {
	UserID      uint64
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	Raiting     float64
	Status      string
}
