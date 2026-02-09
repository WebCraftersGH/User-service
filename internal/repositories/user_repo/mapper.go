package userrepo

import (
	"github.com/WebCraftersGH/User-service/internal/domain"
)

func toGormUser(u domain.User) GormUser {
	return GormUser{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email.String(),
		FIO:           u.FIO,
		BIO:           u.BIO,
		Sex:           u.Sex.String(),
		Birthday:      u.Birthday,
		LastLoginDate: u.LastLoginDate,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func toDomainUser(u GormUser) domain.User {
	return domain.User{
		ID:            u.ID,
		Username:      u.Username,
		Email:         domain.Email(u.Email),
		FIO:           u.FIO,
		BIO:           u.BIO,
		Sex:           domain.NewSexEnum(u.Sex),
		Birthday:      u.Birthday,
		LastLoginDate: u.LastLoginDate,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
