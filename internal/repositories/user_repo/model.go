package userrepo

import (
	"github.com/google/uuid"
	"time"
)

type GormUser struct {
	ID            uuid.UUID
	Username      string
	Email         string
	FIO           string
	BIO           string
	Sex           string
	Birthday      *time.Time
	LastLoginDate *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
