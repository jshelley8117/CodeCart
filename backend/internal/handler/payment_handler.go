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

type PaymentHandler struct {
	PaymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) PaymentHandler {
	return PaymentHandler{
		PaymentService: paymentService,
	}
}

func (ph PaymentHandler) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleCreatePayment")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request model.CreatePaymentRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Warn(common.ERR_REQ_BODY_READ_FAIL, zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Warn(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		z.Warn(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := ph.PaymentService.CreatePayment(r.Context(), request, firebaseUID); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph PaymentHandler) HandleGetPaymentById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleGetPaymentById")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	payment, err := ph.PaymentService.GetPaymentById(r.Context(), id, firebaseUID)
	if err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	paymentApiResponse, err := json.Marshal(payment)
	if err != nil {
		z.Error(common.ERR_REQ_MARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_RETRIEVAL_FAIL, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(paymentApiResponse)
}

func (ph PaymentHandler) HandleUpdatePaymentById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleUpdatePaymentById")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	var request model.UpdatePaymentRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Error(common.ERR_REQ_BODY_READ_FAIL, zap.Error(err))
		http.Error(w, common.ERR_REQ_BODY_READ_FAIL, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		z.Error(common.ERR_REQ_UNMARSH_FAIL, zap.Error(err))
		http.Error(w, common.ERR_REQ_UNMARSH_FAIL, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(request); err != nil {
		z.Error(common.ERR_VALIDATION_FAIL, zap.Error(err))
		http.Error(w, common.ERR_VALIDATION_FAIL, http.StatusBadRequest)
		return
	}

	if err := ph.PaymentService.UpdatePaymentById(r.Context(), request, id, firebaseUID); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_PERSISTENCE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph PaymentHandler) HandleDeletePaymentById(w http.ResponseWriter, r *http.Request) {
	z := utils.FromContext(r.Context(), zap.NewNop())
	z.Debug("Entered HandleDeletePaymentById")

	firebaseUID, ok := r.Context().Value(common.ContextKeyFirebaseUID).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	if err := ph.PaymentService.DeletePaymentById(r.Context(), id, firebaseUID); err != nil {
		z.Error("service invocation failed", zap.Error(err))
		http.Error(w, common.ERR_CLIENT_DB_DELETE_FAIL, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
