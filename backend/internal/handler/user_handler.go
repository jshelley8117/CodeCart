package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
)

var validate = validator.New()

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
	log.Println("Entered HandleCreateUser")
	var request model.CreateUserRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := uh.UserService.CreateUser(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// produce successful 2xx response here - NO BODY
	w.WriteHeader(http.StatusCreated)

}
