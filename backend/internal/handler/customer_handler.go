package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CustomerHandler struct {
	CustomerService service.CustomerService
	Logger          *zap.Logger
}

func NewCustomerHandler(customerService service.CustomerService, logger *zap.Logger) CustomerHandler {
	return CustomerHandler{
		CustomerService: customerService,
		Logger:          logger,
	}
}

func (ch CustomerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), ch.Logger).Named("customer_handler")
	var request model.CreateCustomerRequest
	zLog.Debug("entered HandleCreateCustomer")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zLog.Warn("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		zLog.Warn("go unmarshaling failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		zLog.Warn("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := ch.CustomerService.CreateCustomer(r.Context(), request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
