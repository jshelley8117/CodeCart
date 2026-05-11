package service

import (
	"context"
	"strings"
	"time"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type UserService struct {
	UserPersistence persistence.UserPersistence
	FirebaseAuth    *firebaseauth.Client
}

func NewUserService(userPersistence persistence.UserPersistence, fbAuth *firebaseauth.Client) UserService {
	return UserService{
		UserPersistence: userPersistence,
		FirebaseAuth:    fbAuth,
	}
}

func (us UserService) CreateUser(ctx context.Context, request model.CreateUserRequest) error {
	z := utils.FromContext(ctx, zap.NewNop())
	z.Debug("entered CreateUser")
	userDomainModel := model.User{
		Email:      strings.ToLower(request.Email),
		CustomerId: request.CustomerId,
		AuthId:     request.AuthId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
		Role:       request.Role,
	}

	if err := us.UserPersistence.PersistCreateUser(ctx, userDomainModel); err != nil {
		z.Error("persistence invocation failed: %w", zap.Error(err))
		return err
	}

	if request.Role != nil {
		if err := us.FirebaseAuth.SetCustomUserClaims(ctx, request.AuthId, map[string]any{"role": *request.Role}); err != nil {
			z.Error("failed to set firebase custom claims", zap.Error(err))
			return err
		}
	}

	return nil
}
