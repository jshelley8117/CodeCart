package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"
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
	UserService  service.UserService
	FirebaseAuth *firebaseauth.Client
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		UserService:  userService,
		FirebaseAuth: userService.FirebaseAuth,
	}
}

func (uh UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	var request model.CreateUserRequest
	z.Debug("entered HandleCreateUser")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Warn("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Warn("json deserialization failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&request); err != nil {
		z.Warn("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := uh.UserService.CreateUser(r.Context(), request); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// produce successful 2xx response here - NO BODY
	w.WriteHeader(http.StatusCreated)

}

func (uh UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleRegisterUser")

	// Extract JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		z.Warn("missing or malformed authorization header")
		http.Error(w, "missing or malformed authorization header", http.StatusUnauthorized)
		return
	}
	rawToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify Firebase JWT
	token, err := uh.FirebaseAuth.VerifyIDToken(r.Context(), rawToken)
	if err != nil {
		z.Error("firebase token verification failed", zap.Error(err))
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Read request body
	var registerRequest model.RegisterUserRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Warn("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &registerRequest); err != nil {
		z.Warn("json deserialization failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&registerRequest); err != nil {
		z.Warn("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	// Validate that the AuthId in request matches the verified JWT UID
	if registerRequest.AuthId != token.UID {
		z.Warn("auth id mismatch", zap.String("request_auth_id", registerRequest.AuthId), zap.String("token_uid", token.UID))
		http.Error(w, "auth id does not match token", http.StatusUnauthorized)
		return
	}

	// Default role to "customer" if not provided
	role := registerRequest.Role
	if role == nil {
		defaultRole := "customer"
		role = &defaultRole
	}

	// Create CreateUserRequest with verified Firebase UID
	createUserRequest := model.CreateUserRequest{
		Email:  registerRequest.Email,
		AuthId: token.UID,
		Role:   role,
	}

	// Call service to create user
	if err := uh.UserService.CreateUser(r.Context(), createUserRequest); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{
		"message": "user registered successfully",
		"uid":     token.UID,
	}

	data, err := json.Marshal(response)
	if err != nil {
		z.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
