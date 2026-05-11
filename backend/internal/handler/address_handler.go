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

type AddressHandler struct {
	AddressService service.AddressService
}

func NewAddressHandler(addressService service.AddressService) AddressHandler {
	return AddressHandler{
		AddressService: addressService,
	}
}

func (ah AddressHandler) HandleCreateAddress(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleCreateAddress")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	if err := ah.AddressService.CreateAddress(r.Context(), request, firebaseUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ah AddressHandler) HandleGetAllAddressesById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleGetAllAddresses")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	addresses, err := ah.AddressService.GetAddressesByAuthId(r.Context(), firebaseUID)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addressesApiResponse, err := json.Marshal(addresses)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(addressesApiResponse)
}

func (ah AddressHandler) HandleGetAddressById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleGetAddressById")

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

	address, err := ah.AddressService.GetAddressById(r.Context(), id)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addressApiResponse, err := json.Marshal(address)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(addressApiResponse)

}

func (ah AddressHandler) HandleUpdateAddressById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleUpdateAddressById")

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

	var request model.UpdateAddressRequest

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

	err = ah.AddressService.UpdateAddressById(r.Context(), request, id)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (ah AddressHandler) HandleDeleteAddressById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("entered HandleDeleteAddressById")

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

	err = ah.AddressService.DeleteAddressById(r.Context(), id)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
