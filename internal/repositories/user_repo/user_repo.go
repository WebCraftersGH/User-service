package userrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/WebCraftersGH/User-service/internal/domain"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo() usecase.UserRepo {
	return &userRepo{}
}

func (r *userRepo) Create(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {

	gormUser := toGormUser(u)

	err := r.db.WithContext(ctx).Create(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserAlreadyExists, err)
		}
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	user := toDomainUser(gormUser)

	return user, nil
}

func (r *userRepo) Read(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {

	var gormUser GormUser

	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	domainUser := toDomainUser(gormUser)

	return domainUser, nil
}

func (r *userRepo) Update(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {

	gormUser := toGormUser(u)

	result := r.db.WithContext(ctx).Where("id = ?", u.ID).Updates(gormUser)

	if result.Error != nil {
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.User{}, domain.ErrUserNotFound
	}

	var updatedGormUser GormUser
	err := r.db.WithContext(ctx).Where("id = ?", u.ID).First(&updatedGormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	domainUser := toDomainUser(updatedGormUser)

	return domainUser, nil
}

func (r *userRepo) Delete(
	ctx context.Context,
	userID uuid.UUID,
) error {

	var gormUser GormUser

	err := r.db.WithContext(ctx).Where("id = ?", userID).Delete(gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		return fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	return nil
}
