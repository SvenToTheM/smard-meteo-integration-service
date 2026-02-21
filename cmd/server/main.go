package main

import (
	"log/slog"
	"net/http"

	httpHandler "github.com/sventothem/smard-meteo-integration-service/internal/adapters/httphandler"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/openmeteo"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/smard"
	"github.com/sventothem/smard-meteo-integration-service/internal/core/services"
)

func main() {
	logger := slog.Default()
	
	sClient := smard.NewClient()
	mClient := openmeteo.NewClient()

	analysisService := services.NewGridService(sClient, mClient)

	handler := httpHandler.NewAnalysisHandler(analysisService)

	mux := http.NewServeMux()

	// register Route
	mux.HandleFunc("/api/v1/analysis", handler.GetAnalysis)

	logger.Info("Service gestartet", "port", 8080)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}