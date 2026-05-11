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

type InventoryHandler struct {
	InventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) InventoryHandler {
	return InventoryHandler{
		InventoryService: inventoryService,
	}
}

func (ih InventoryHandler) HandleCreateInventory(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleCreateInventory")

	var request model.CreateInventoryRequest

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

	if err := validate.Struct(request); err != nil {
		z.Error(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := ih.InventoryService.CreateInventory(r.Context(), request); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ih InventoryHandler) HandleGetAllInventory(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleGetAllInventory")

	page, pageSize, err := utils.ParsePaginationInput(r.Context(), r)
	if err != nil {
		z.Error("failed to parse pagination input", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	inventory, total, err := ih.InventoryService.GetAllInventory(r.Context(), page, pageSize)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	totalPages := utils.CalculateTotalPages(int(total), pageSize)

	response := common.PaginatedResponse{
		Data:       inventory,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	inventoryApiResponse, err := json.Marshal(response)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(inventoryApiResponse)
}

func (ih InventoryHandler) HandleGetInventoryById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleGetInventoryById")

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

	item, err := ih.InventoryService.GetInventoryById(r.Context(), id)
	if err != nil {
		z.Error("service invocation failed", zap.Int("id", id), zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	inventoryApiResponse, err := json.Marshal(item)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(inventoryApiResponse)
}

func (ih InventoryHandler) HandleUpdateInventoryById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleUpdateInventoryById")

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

	var request model.UpdateInventoryRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Error(common.ERR_REQ_BODY_READ_FAIL, zap.Int("id", id), zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Error(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := ih.InventoryService.UpdateInventoryById(r.Context(), id, request); err != nil {
		z.Error("service invocation failed", zap.Int("id", id), zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ih InventoryHandler) HandleDeleteInventoryById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleDeleteInventoryById")

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

	if err := ih.InventoryService.DeleteInventoryById(r.Context(), id); err != nil {
		z.Error("service invocation failed", zap.Int("id", id), zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_DELETE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
