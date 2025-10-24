package models

import "time"

type User struct {
	ID                string
	Username          string
	Email             string
	PasswordHash      string
	VerificationToken string
	Role              string
	CreatedAt         time.Time
}
