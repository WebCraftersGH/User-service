package controller

import (
	"github.com/WebCraftersGH/User-service/internal/domain"
)

func toUserResponse(u domain.User) UserResponse {
	return UserResponse{
		ID:       u.ID.String(),
		Username: u.Username,
		Email:    u.Email.String(),
		FIO:      u.FIO,
		BIO:      u.BIO,
		Sex:      u.Sex.String(),
		Birthday: u.Birthday,
	}
}
