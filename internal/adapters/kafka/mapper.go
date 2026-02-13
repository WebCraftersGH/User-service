package kafka

import (
	"github.com/WebCraftersGH/User-service/internal/domain"
)

func toDomainUser(u KafkaUser) domain.User {
	return domain.User{
		Username: u.Username,
		Email:    domain.Email(u.Email),
		FIO:      u.FIO,
		BIO:      u.BIO,
		Sex:      domain.NewSexEnum(u.Sex),
		Birthday: u.Birthday,
	}
}
