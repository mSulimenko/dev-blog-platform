package dto

import "encoding/json"

type UserRegisteredEvent struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

func (e *UserRegisteredEvent) ToJson() ([]byte, error) {
	return json.Marshal(e)
}
