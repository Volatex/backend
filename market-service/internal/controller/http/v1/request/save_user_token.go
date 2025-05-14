package request

type SaveUserToken struct {
	TinkoffToken string `json:"tinkoff_token" validate:"required"`
}
