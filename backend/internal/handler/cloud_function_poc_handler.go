package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jshelley8117/CodeCart/internal/service"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

type CloudFunctionHandler struct {
	CloudFunctionService service.CloudFunctionService
}

func NewCloudFunctionHandler(cloudFunctionService service.CloudFunctionService) CloudFunctionHandler {
	return CloudFunctionHandler{
		CloudFunctionService: cloudFunctionService,
	}
}

func (cfh CloudFunctionHandler) HandleGetHelloWorld(w http.ResponseWriter, r *http.Request) {
	zLog := utils.FromContext(r.Context(), zap.NewNop())
	zLog.Debug("entered HandleGetHelloWorld")

	response, err := cfh.CloudFunctionService.GetHelloWorld(r.Context())
	if err != nil {
		zLog.Error("service invocation failed", zap.Error(err))
		http.Error(w, "Failed to invoke cloud function", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		zLog.Error("go marshaling failed", zap.Error(err))
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
