package controller

import (
	"time"
)

type UserResponse struct {
	ID            string
	Username      string
	Email         string
	FIO           string
	BIO           string
	Sex           string
	Birthday      *time.Time
	LastLoginDate *time.Time
}

type UserUpdateRequest struct {
	Username string
	FIO string
	BIO string
	Sex string
	Birthday *time.Time
}
