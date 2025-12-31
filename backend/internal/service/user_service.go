package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/persistence"
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
	log.Println("Entered CreateUser")
	userDomainModel := model.User{
		Email:      strings.ToLower(request.Email),
		CustomerId: request.CustomerId,
		GCAuthId:   request.GCAuthId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
	}

	if err := us.UserPersistence.PersistCreateUser(ctx, userDomainModel); err != nil {
		// log error here
		return err
	}

	return nil
}
