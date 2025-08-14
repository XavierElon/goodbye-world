package domain

import (
	"time"
)

type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64 `json:"expires_in"`
	ExpiresAt time.Time `json:"expires_at"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthResponse struct {
	User *User `json:"user"`
	Token *AuthToken `json:"token"`
}

type AuthError struct {
	Error string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// What goes inside the JWT token
type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Exp         int64  `json:"exp"`
}