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

type OrderDiscountHandler struct {
	OrderDiscountService service.OrderDiscountService
}

func NewOrderDiscountHandler(orderDiscountService service.OrderDiscountService) OrderDiscountHandler {
	return OrderDiscountHandler{
		OrderDiscountService: orderDiscountService,
	}
}

func (odh OrderDiscountHandler) HandleCreateOrderDiscount(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleCreateOrderDiscount")

	orderIdPathVal := r.PathValue("orderId")
	if orderIdPathVal == "" {
		z.Error("orderId field in endpoint path parameter is missing")
		http.Error(w, "orderId is empty", http.StatusBadRequest)
		return
	}

	orderId, err := strconv.Atoi(orderIdPathVal)
	if err != nil {
		z.Error("failed to convert orderId value from string to integer")
		http.Error(w, "server failed to process orderId value", http.StatusInternalServerError)
		return
	}

	var request model.OrderDiscountRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Error(common.ERR_REQ_BODY_READ_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Error(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	request.OrderId = orderId

	if err := validate.Struct(request); err != nil {
		z.Error(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := odh.OrderDiscountService.CreateOrderDiscount(r.Context(), request); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (odh OrderDiscountHandler) HandleGetOrderDiscountById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleGetOrderDiscountById")

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

	orderDiscount, err := odh.OrderDiscountService.GetOrderDiscountById(r.Context(), id)
	if err != nil {
		z.Error("service invocation failed", zap.Int("id", id), zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	apiResponse, err := json.Marshal(orderDiscount)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(apiResponse)
}
