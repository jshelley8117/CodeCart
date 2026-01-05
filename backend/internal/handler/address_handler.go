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

type AddressHandler struct {
	AddressService service.AddressService
}

func NewAddressHandler(addressService service.AddressService) AddressHandler {
	return AddressHandler{
		AddressService: addressService,
	}
}

func (ah AddressHandler) HandleCreateAddress(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered HandleCreateAddress")

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
