package entity

type User struct {
	Id              int64  `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	IsEmailVerified bool   `json:"is_email_verified"`
}
