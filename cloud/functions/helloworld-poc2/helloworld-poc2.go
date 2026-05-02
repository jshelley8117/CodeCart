package helloworldpoc2

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorldPOC2", HelloWorldPOC2)
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func HelloWorldPOC2(w http.ResponseWriter, r *http.Request) {
	response := HelloWorldResponse{
		Message: "Hello World from Cloud Function 2!",
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
