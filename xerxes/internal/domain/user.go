package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
	IsVerified bool      `json:"is_verified"`
}

type UserRegistration struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UserLogin struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Code        string `json:"code" validate:"required"`
}

type VerificationRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerificationResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}