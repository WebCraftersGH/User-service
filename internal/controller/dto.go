package controller

import (
	"time"
)

type UserResponse struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	FIO      string     `json:"fio"`
	BIO      string     `json:"bio"`
	Sex      string     `json:"sex"`
	Birthday *time.Time `json:"birthday,omitempty"`
}

type UserUpdateRequest struct {
	Username string     `json:"username"`
	FIO      string     `json:"fio"`
	BIO      string     `json:"bio"`
	Sex      string     `json:"sex"`
	Birthday *time.Time `json:"birthday,omitempty"`
}
