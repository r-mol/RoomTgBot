package internal

type User struct {
	ID int64 `json:"id"`

	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}
