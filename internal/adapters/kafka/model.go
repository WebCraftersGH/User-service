package kafka

import (
	"time"
)

type KafkaUser struct {
	Username string
	Email    string
	FIO      string `json:"fio"`
	BIO      string `json:"bio"`
	Sex      string
	Birthday *time.Time
}
