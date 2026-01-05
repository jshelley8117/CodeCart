package service

import (
	"context"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type UserService struct {
	UserPersistence persistence.UserPersistence
	Logger          *zap.Logger
}

func NewUserService(userPersistence persistence.UserPersistence, logger *zap.Logger) UserService {
	return UserService{
		UserPersistence: userPersistence,
		Logger:          logger,
	}
}

func (us UserService) CreateUser(ctx context.Context, request model.CreateUserRequest) error {
	zLog := utils.FromContext(ctx, us.Logger).Named("user_service")
	zLog.Debug("entered CreateUser")
	userDomainModel := model.User{
		Email:      strings.ToLower(request.Email),
		CustomerId: request.CustomerId,
		GCAuthId:   request.GCAuthId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
	}

	if err := us.UserPersistence.PersistCreateUser(ctx, userDomainModel); err != nil {
		zLog.Error("persistence invocation failed: %w", zap.Error(err))
		return err
	}

	return nil
}
