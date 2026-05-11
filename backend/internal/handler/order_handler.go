package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type OrderHandler struct {
	OrderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return OrderHandler{
		OrderService: orderService,
	}
}

func (oh OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())

	var request model.CreateOrderRequest

	z.Debug("Entered HandleCreateOrder")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Warn(common.ERR_REQ_BODY_READ_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Warn(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		z.Warn(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}
	if err := oh.OrderService.CreateOrder(r.Context(), request); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (oh OrderHandler) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleGetAllOrders")

	orders, err := oh.OrderService.GetAllOrders(r.Context())
	if err != nil {
		z.Error("Service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	ordersApiResponse, err := json.Marshal(orders)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersApiResponse)
}

func (oh OrderHandler) HandleFetchOrderById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleFetchOrderById")

	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		z.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idPathVal)
	if err != nil {
		z.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	orders, err := oh.OrderService.FetchOrderById(r.Context(), id)
	if err != nil {
		z.Error("Service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	ordersApiResponse, err := json.Marshal(orders)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersApiResponse)

}

func (oh OrderHandler) HandleUpdateOrderById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandlePersistUpdateOrderById")

	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		z.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idPathVal)
	if err != nil {
		z.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	var request model.UpdateOrderRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Error("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Error("go unmarshaling failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		z.Error("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := oh.OrderService.UpdateOrderById(r.Context(), request, id); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
