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

type DiscountHandler struct {
	DiscountService service.DiscountService
}

func NewDiscountHandler(discountService service.DiscountService) DiscountHandler {
	return DiscountHandler{
		DiscountService: discountService,
	}
}

func (dh DiscountHandler) HandleCreateDiscount(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleCreateDiscount")

	var request model.CreateDiscount
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

	if err := dh.DiscountService.CreateDiscount(r.Context(), request); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
