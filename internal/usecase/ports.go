package usecase

import (
	"context"
	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetAllUser(ctx context.Context, limit, offset int) ([]domain.User, error)

	CreateUser(ctx context.Context, u domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, u domain.User) (domain.User, error)

	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UserRepo interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	Read(ctx context.Context, userID uuid.UUID) (domain.User, error)
	Update(ctx context.Context, u domain.User) (domain.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type Consumer interface {
	Start()
	Stop() error
}
