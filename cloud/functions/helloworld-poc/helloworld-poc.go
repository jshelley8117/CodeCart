package helloworldpoc

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorldPOC", HelloWorldPOC)
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func HelloWorldPOC(w http.ResponseWriter, r *http.Request) {
	response := HelloWorldResponse{
		Message: "Hello World from Cloud Function!",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
