package dto

type UserRegisteredEvent struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
}
