package usecase

import (
	"context"
	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/google/uuid"
)

type userService struct {
	repo UserRepo
	lg   Logger
}

var _ UserService = (*userService)(nil)

func NewUserService(repo UserRepo, lg Logger) *userService {
	return &userService{
		repo: repo,
		lg:   lg,
	}
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	s.lg.Debug("[User-service][GetUser][DEBUG] - Start function")

	user, err := s.repo.Read(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	s.lg.Debug("[User-service][GetUser][DEBUG] - End function")
	return user, nil
}

func (s *userService) GetAllUser(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return nil, nil
}

func (s *userService) CreateUser(ctx context.Context, u domain.User) (domain.User, error) {
	s.lg.Debug("[User-service][CreateUser][DEBUG] - Start function")

	user, err := s.repo.Create(ctx, u)
	if err != nil {
		return domain.User{}, err
	}

	s.lg.Debug("[User-service][CreateUser][DEBUG] - End function")
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, u domain.User) (domain.User, error) {
	s.lg.Debug("[User-service][UpdateUser][DEBUG] - Start function")

	user, err := s.repo.Update(ctx, u)
	if err != nil {
		return domain.User{}, err
	}

	s.lg.Debug("[User-service][UpdateUser][DEBUG] - End function")
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	s.lg.Debug("[User-service][DeleteUser][DEBUG] - Start function")

	if err := s.repo.Delete(ctx, userID); err != nil {
		return err
	}

	s.lg.Debug("[User-service][DeleteUser][DEBUG] - End function")
	return nil
}
