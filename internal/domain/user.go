package domain

import (
	"github.com/google/uuid"
	"time"
)

type SexEnum int

func NewSexEnum(value string) SexEnum {
	switch value {
	case "Male":
		return SexEnum(0)
	case "Female":
		return SexEnum(1)
	default:
		return SexEnum(2)
	}
}

const (
	SexMale SexEnum = iota
	SexFemale
	SexOther
)

var SexEnumName = map[SexEnum]string{
	SexMale:   "Male",
	SexFemale: "Female",
	SexOther:  "Other",
}

func (e SexEnum) String() string {
	return SexEnumName[e]
}

type Email string

func (e Email) String() string {
	return string(e)
}

type User struct {
	ID            uuid.UUID
	Username      string
	Email         Email
	FIO           string
	BIO           string
	Sex           SexEnum
	Birthday      *time.Time
	LastLoginDate *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
