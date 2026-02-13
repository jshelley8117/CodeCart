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

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return ProductHandler{
		ProductService: productService,
	}
}

func (ph ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleCreateProduct")

	var request model.CreateProductRequest

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

	// ERROR: Request fails here
	if err := validate.Struct(request); err != nil {
		zLog.Warn(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusBadRequest)
		return
	}

	if err := ph.ProductService.ServiceCreateProduct(r.Context(), request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph ProductHandler) HandleFetchAllProducts(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleFetchAllProducts")

	page := 1
	pageSize := 10

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	products, total, err := ph.ProductService.ServiceFetchAllProducts(r.Context(), page, pageSize)
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response := common.PaginatedResponse{
		Data:       products,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	productsApiResponse, err := json.Marshal(response)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productsApiResponse)
}

func (ph ProductHandler) HandleFetchProductById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleFetchProductById")

	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		zLog.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}

	productId, err := strconv.Atoi(idPathVal)
	if err != nil {
		zLog.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	product, err := ph.ProductService.ServiceFetchProductById(r.Context(), productId)
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	productApiResponse, err := json.Marshal(product)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productApiResponse)
}

func (ph ProductHandler) HandleUpdateProductById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleUpdateProductById")

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

	var request model.UpdateProductRequest

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

	if err := ph.ProductService.ServiceUpdateProductById(r.Context(), id, request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph ProductHandler) HandleFetchAllProductVariantsByProductId(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleFetchAllProductVariantsByProductId")

	idPathVal := r.PathValue("id")
	if idPathVal == "" {
		zLog.Error("ID field in endpoint path parameter is missing")
		http.Error(w, "ID is empty", http.StatusBadRequest)
		return
	}

	productId, err := strconv.Atoi(idPathVal)
	if err != nil {
		zLog.Error("failed to convert id value from string to integer")
		http.Error(w, "server failed to process ID value", http.StatusInternalServerError)
		return
	}

	page := 1
	pageSize := 10

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	variants, total, err := ph.ProductService.ServiceFetchAllProductVariantsByProductId(r.Context(), productId, page, pageSize)
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response := common.PaginatedResponse{
		Data:       variants,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	variantsApiResponse, err := json.Marshal(response)
	if err != nil {
		zLog.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(variantsApiResponse)
}

func (ph ProductHandler) HandleUpdateProductVariantById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleUpdateProductVariantById")

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

	var request model.UpdateProductVariant

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

	if err := ph.ProductService.ServiceUpdateProductVariantById(r.Context(), id, request); err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph ProductHandler) HandleDeleteProductById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleDeleteProductById")

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

	err = ph.ProductService.ServiceDeleteProductById(r.Context(), id)
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (ph ProductHandler) HandleDeleteProductVariantById(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleDeleteProductVariantById")

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

	err = ph.ProductService.ServiceDeleteProductVariantById(r.Context(), id)
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_REQUEST_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
