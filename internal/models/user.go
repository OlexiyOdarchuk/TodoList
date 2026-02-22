package models

import "time"

type User struct {
	Id                      string    `json:"id" db:"id"`
	Username                string    `json:"username" db:"username"`
	Email                   string    `json:"email" db:"email"`
	PendingEmail            string    `json:"pending_email" db:"pending_email"`
	PasswordHash            string    `json:"-" db:"password_hash"`
	IsVerified              bool      `json:"is_verified" db:"is_verified"`
	VerificationCode        string    `json:"-" db:"verification_code"`
	VerificationCodeExpires time.Time `json:"-" db:"verification_code_expires"`
	OauthProvider           string    `json:"oauth_provider" db:"oauth_provider"`
	OauthId                 string    `json:"oauth_id" db:"oauth_id"`
}
