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
	Logger       *zap.Logger
}

func NewOrderHandler(orderService service.OrderService, logger *zap.Logger) OrderHandler {
	return OrderHandler{
		OrderService: orderService,
		Logger:       logger,
	}
}

func (oh OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), oh.Logger).Named("order_handler")

	var request model.CreateOrderRequest

	zLog.Debug("Entered HandleCreateOrder")

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
	if err := oh.OrderService.CreateOrder(r.Context(), request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (oh OrderHandler) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), oh.Logger).Named("order_handler")
	zLog.Debug("entered HandleGetAllOrders")

	orders, err := oh.OrderService.GetAllOrders(r.Context())
	if err != nil {
		zLog.Error("Service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	ordersApiResponse, err := json.Marshal(orders)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersApiResponse)
}

func (oh OrderHandler) HandleFetchOrderById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), oh.Logger).Named("order_handler")
	zLog.Debug("entered HandleFetchOrderById")

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

	orders, err := oh.OrderService.FetchOrderById(r.Context(), id)
	if err != nil {
		zLog.Error("Service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	ordersApiResponse, err := json.Marshal(orders)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersApiResponse)

}

func (oh OrderHandler) HandleUpdateOrderById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), oh.Logger).Named("order_handler")
	zLog.Debug("entered HandlePersistUpdateOrderById")

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

	var request model.UpdateOrderRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zLog.Error("request body read failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		zLog.Error("go unmarshaling failed", zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		zLog.Error("struct validation failed", zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := oh.OrderService.UpdateOrderById(r.Context(), request, id); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
