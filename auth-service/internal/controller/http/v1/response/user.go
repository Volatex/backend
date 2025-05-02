package response

type User struct {
	Id    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
}
