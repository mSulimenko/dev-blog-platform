package dto

import "time"

type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=50"`
}

type UserUpdateRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=30"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=4,max=50"`
	Role     *string `json:"role" validate:"oneof=unverified user admin"`
}

type UserResp struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
