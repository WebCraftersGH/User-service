package userrepo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GormUser struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username      string
	Email         string
	FIO           string
	BIO           string
	Sex           string
	Birthday      *time.Time
	LastLoginDate *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (GormUser) TableName() string {
	return "users"
}
