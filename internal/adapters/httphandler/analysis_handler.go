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
	return &AnalysisHandler{
		service: service,
	}
}

func (h *AnalysisHandler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	// Optional http path param for retrieving also the data where generation is not covering consumption
	all := r.URL.Query().Get("all") == "false"

	slog.Info("Received analysis request",
		slog.String("start_param", startStr),
		slog.String("end_param", endStr),
		slog.Bool("all", all),
		slog.String("client_ip", r.RemoteAddr),
	)

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		slog.Warn("Invalid start date", slog.String("input", startStr), slog.Any("error", err))
		http.Error(w, "Invalid start date format (RFC3339 required)", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		slog.Warn("Invalid end date", slog.String("input", endStr), slog.Any("error", err))
		http.Error(w, "Invalid end date format (RFC3339 required)", http.StatusBadRequest)
		return
	}

	results, err := h.service.Analyze(r.Context(), start, end, all)
	if err != nil {
		slog.Error("Service analysis failed", slog.Time("start", start), slog.Time("end", end), slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Successfully processed analysis request", slog.Int("results_count", len(results)))

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
