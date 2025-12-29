package handler

import (
	"context"
	"encoding/json"
	"io"
	"json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
)

var validate validator.Validate

type UserService interface {
	CreateUser(ctx context.Context, request model.CreateUserRequest) error
}

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (uh UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request model.CreateUserRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return
	}

	if err := validate.Struct(request); err != nil {
		return
	}

	if err := uh.UserService.CreateUser(r.Context(), request); err != nil {
		return
	}

	// produce successful 2xx response here

}
