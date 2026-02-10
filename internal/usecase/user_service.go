package usecase

import (
	"context"
	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/google/uuid"
)

type userService struct {
	repo UserRepo
}

var _ UserService = (*userService)(nil)

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return s.repo.Read(ctx, userID)
}

func (s *userService) GetAllUser(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return nil, nil
}

func (s *userService) CreateUser(ctx context.Context, u domain.User) (domain.User, error) {
	return s.repo.Create(ctx, u)
}

func (s *userService) UpdateUser(ctx context.Context, u domain.User) (domain.User, error) {
	return s.repo.Update(ctx, u)
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Delete(ctx, userID)
}
