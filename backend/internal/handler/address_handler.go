package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type AddressHandler struct {
	AddressService service.AddressService
	Logger         *zap.Logger
}

func NewAddressHandler(addressService service.AddressService) AddressHandler {
	return AddressHandler{
		AddressService: addressService,
	}
}

func (ah AddressHandler) HandleCreateAddress(w http.ResponseWriter, r *http.Request) {
	zLog := ah.getZLog(r.Context())
	zLog.Debug("Entered HandleCreateAddress")

	var request model.CreateAddressRequest

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
		return
	}

	if err := ah.AddressService.CreateAddress(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ah AddressHandler) HandleGetAllAddresses(w http.ResponseWriter, r *http.Request) {
	zLog := ah.getZLog(r.Context())
	zLog.Debug("Entered HandleGetAllAddresses")

	addresses, err := ah.AddressService.GetAllAddresses(r.Context())
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addressesApiResponse, err := json.Marshal(addresses)
	if err != nil {
		zLog.Error("go marshaling failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(addressesApiResponse)
}

func (ah AddressHandler) getZLog(ctx context.Context) *zap.Logger {
	return utils.FromContext(ctx, ah.Logger)
}
