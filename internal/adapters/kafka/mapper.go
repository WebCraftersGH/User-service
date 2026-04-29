package kafka

import (
	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/google/uuid"
)

func toDomainUser(u KafkaUser) domain.User {
	uID, _ := uuid.Parse(u.UserID)
	return domain.User{
		Email:    domain.Email(u.Email),
		CreatedAt: u.CreatedAt,
		ID: uID,
	}
}
