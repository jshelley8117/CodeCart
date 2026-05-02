package handler

import (
	"encoding/json"
	"io"
	"net/http"

	// "strconv"

	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/model"
	"github.com/jshelley8117/CodeCart/internal/service"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	InventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) InventoryHandler {
	return InventoryHandler{
		InventoryService: inventoryService,
	}
}

func (ih InventoryHandler) HandleCreateInventory(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("Entered HandleCreateInventory")

	var request model.CreateInventoryRequest

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

	if err := ih.InventoryService.CreateInventory(r.Context(), request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (ih InventoryHandler) HandleGetAllInventory(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("Entered HandleGetAllInventory")

	inventory, err := ih.InventoryService.GetAllInventory(r.Context())
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inventoryApiResponse, err := json.Marshal(inventory)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(inventoryApiResponse)
}

// func (ih InventoryHandler) HandleGetInventoryById(w http.ResponseWriter, r *http.Request) {
// 	zLog := utils.FromContext(r.Context(), zap.NewNop())
// 	zLog.Debug("Entered HandleGetInventoryById")

// 	idPathVal := r.PathValue("id")
// 	if idPathVal == "" {
// 		zLog.Error("ID field in endpoint path parameter is missing")
// 		http.Error(w, "ID is empty", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := strconv.Atoi(idPathVal)
// 	if err != nil {
// 		zLog.Error("failed to convert id value from string to integer")
// 		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
// 		return
// 	}

// 	var request model.Inventory

// 	body, err := ih.InventoryService.GetAllInventory(r.Context(), id)
// }

func (ih InventoryHandler) HandleUpdateInventoryById(w http.ResponseWriter, r *http.Request)

func (ih InventoryHandler) HandleDeleteInventoryById(w http.ResponseWriter, r *http.Request)
