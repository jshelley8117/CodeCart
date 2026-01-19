package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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
		Logger:          logger.Named("customer_handler"),
	}
}

func (ch CustomerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	zLog := ch.getZLog(r.Context())
	var request model.CreateCustomerRequest
	zLog.Debug("entered HandleCreateCustomer")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zLog.Warn(common.ERR_REQ_BODY_READ_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		zLog.Warn(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		zLog.Warn(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := ch.CustomerService.CreateCustomer(r.Context(), request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (ch CustomerHandler) HandleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	zLog := ch.getZLog(r.Context())
	zLog.Debug("entered HandleGetAllCustomers")

	customers, err := ch.CustomerService.GetAllCustomers(r.Context())
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, "Failed to retrieve customers", http.StatusInternalServerError)
		return
	}

	customersApiResponse, err := json.Marshal(customers)
	if err != nil {
		zLog.Error("go marshaling failed", zap.Error(err))
		http.Error(w, "Failed to serialize response to client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(customersApiResponse)
}

func (ch CustomerHandler) HandleDeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	zLog := ch.getZLog(r.Context())
	zLog.Debug("entered HandleDeleteCustomerById")
	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		zLog.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idPathVal)
	if err != nil {
		zLog.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	if err := ch.CustomerService.DeleteCustomerById(r.Context(), id); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, fmt.Sprintf("Failed to delete customer [ID: %v]", id), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ch CustomerHandler) HandleUpdateCustomerById(w http.ResponseWriter, r *http.Request) {
	zLog := ch.getZLog(r.Context())
	zLog.Debug("entered HandleUpdateCustomerById")

	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		zLog.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idPathVal)
	if err != nil {
		zLog.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	var request model.UpdateCustomerRequest

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

	if err := ch.CustomerService.UpdateCustomerById(r.Context(), request, id); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, fmt.Sprintf("Failed to update customer [ID: %d]: %v", id, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ch CustomerHandler) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, ch.Logger)
}
