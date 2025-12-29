package service

import (
	"context"
	"strings"
	"time"

	"github.com/jshelley8117/CodeCart/internal/model"
)

type UserService struct {
}

func NewUserService() UserService {
	return UserService{}
}

func (us UserService) CreateUser(ctx context.Context, request model.CreateUserRequest) error {
	userDomainModel := model.User{
		Email:     strings.ToLower(request.Email),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

}
