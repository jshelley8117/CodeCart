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
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

var validate = validator.New()

type UserService interface {
	CreateUser(ctx context.Context, request model.CreateUserRequest) error
}

type UserHandler struct {
	UserService service.UserService
	Logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) UserHandler {
	return UserHandler{
		UserService: userService,
		Logger:      logger,
	}
}

func (uh UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), uh.Logger).Named("user_handler")
	var request model.CreateUserRequest
	zLog.Debug("entered HandleCreateUser")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zLog.Warn("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		zLog.Warn("json deserialization failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		zLog.Warn("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := uh.UserService.CreateUser(r.Context(), request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// produce successful 2xx response here - NO BODY
	w.WriteHeader(http.StatusCreated)

}
