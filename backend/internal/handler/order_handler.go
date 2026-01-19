package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
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
	log.Println("Entered HandleCreateOrder")

	var request model.CreateOrderRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		log.Println("Validation failed:", err)
		return
	}
	if err := oh.OrderService.CreateOrder(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
