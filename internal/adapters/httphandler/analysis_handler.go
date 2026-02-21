package httphandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/ports"
)

type AnalysisHandler struct {
	// Only knows the interface
	service ports.AnalysisService
}

// Constructor for handler (DI)
func NewAnalysisHandler(service ports.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler {
		service: service,
	}
}

func (h *AnalysisHandler) GetAnalysis(w http.ResponseWriter, r *http.Request){
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start date format (RFC3339 required)", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end date format (RFC3339 required)", http.StatusBadRequest)
		return
	}
	// Call service with r.Context() for timeouts)
	results, err := h.service.Analyze(r.Context(), start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SmardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "pong"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode json", "error", err)
	}
}
