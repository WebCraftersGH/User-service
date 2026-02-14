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
	lg usecase.Logger
}

func NewUserRepo(db *gorm.DB, lg usecase.Logger) usecase.UserRepo {
	return &userRepo{
		db: db,
		lg: lg,
	}
}

func (r *userRepo) Create(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {
	r.lg.Debug("[User-repo][Create][DEBUG] - Start function", "email", u.Email.String())

	gormUser := toGormUser(u)

	err := r.db.WithContext(ctx).Create(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.lg.Warn("[User-repo][Create][WARN] - Dublicate gorm error", "gorm_err", err)
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserAlreadyExists, err)
		}
		r.lg.Error("[User-repo][Create][ERROR] - Internal db error", "gorm_err", err)
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	user := toDomainUser(gormUser)

	r.lg.Debug("[User-repo][Create][DEBUG] - End function")
	return user, nil
}

func (r *userRepo) Read(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	r.lg.Debug("[User-repo][Read][DEBUG] - Start function", "user_id", userID)

	var gormUser GormUser

	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.lg.Info("[User-repo][Read][INFO] - User not found", "user_id", userID)
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		r.lg.Error("[User-repo][Read][ERROR] - Internal db error", "gorm_err", err)
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	domainUser := toDomainUser(gormUser)

	r.lg.Debug("[User-repo][Read][DEBUG] - End function")
	return domainUser, nil
}

func (r *userRepo) Update(
	ctx context.Context,
	u domain.User,
) (domain.User, error) {
	r.lg.Debug("[User-repo][Update][DEBUG] - Start function", "user_id", u.ID)

	gormUser := toGormUser(u)

	result := r.db.WithContext(ctx).Where("id = ?", u.ID).Updates(gormUser)

	if result.Error != nil {
		r.lg.Error("[User-repo][Update][ERROR] - Internal db error", "gorm_err", result.Error)
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, result.Error)
	}

	if result.RowsAffected == 0 {
		r.lg.Warn("[User-repo][Update][WARN] - User not found", "user_id", u.ID)
		return domain.User{}, domain.ErrUserNotFound
	}

	var updatedGormUser GormUser
	err := r.db.WithContext(ctx).Where("id = ?", u.ID).First(&updatedGormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.lg.Error("[User-repo][Update][Error] - User not found after update", "gorm_err", err, "user_id", u.ID)
			return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		r.lg.Error("[User-repo][Update][ERROR] - Internal db error", "gorm_err", err)
		return domain.User{}, fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	domainUser := toDomainUser(updatedGormUser)

	r.lg.Debug("[User-repo][Update][DEBUG] - End function")
	return domainUser, nil
}

func (r *userRepo) Delete(
	ctx context.Context,
	userID uuid.UUID,
) error {
	r.lg.Debug("[User-repo][Delete][DEBUG] - Start function", "user_id", userID)

	var gormUser GormUser

	err := r.db.WithContext(ctx).Where("id = ?", userID).Delete(gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.lg.Info("[User-repo][Delete][INFO] - User not found")
			return fmt.Errorf("%w: gorm error: %v", domain.ErrUserNotFound, err)
		}
		r.lg.Error("[User-repo][Delete][ERROR] - Internal db error", "gorm_err", err)
		return fmt.Errorf("%w: gorm error: %v", domain.InternalError, err)
	}

	r.lg.Debug("[User-repo][Delete][DEBUG] - End function")
	return nil
}
