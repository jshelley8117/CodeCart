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
}

func NewUserService(userPersistence persistence.UserPersistence) UserService {
	return UserService{
		UserPersistence: userPersistence,
	}
}

func (us UserService) CreateUser(ctx context.Context, request model.CreateUserRequest) error {
	zLog := utils.FromContext(ctx, zap.NewNop())
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
