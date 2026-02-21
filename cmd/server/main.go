package main

import (
	"log/slog"
	"net/http"

	httpHandler "github.com/sventothem/smard-meteo-integration-service/internal/adapters/httphandler"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/openmeteo"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/smard"
	"github.com/sventothem/smard-meteo-integration-service/internal/config"
	"github.com/sventothem/smard-meteo-integration-service/internal/core/services"
)

func main() {
	logger := slog.Default()

	cfg := config.LoadConfig()

	extHttpClient := &http.Client{}

	sClient := smard.NewClient(cfg.SmardBaseURL, extHttpClient)
	mClient := openmeteo.NewClient(cfg.MeteoBaseURL, extHttpClient)

	analysisService := services.NewGridService(sClient, mClient, cfg.DefaultLat, cfg.DefaultLon)

	handler := httpHandler.NewAnalysisHandler(analysisService)

	mux := http.NewServeMux()

	// register Route
	mux.HandleFunc("/api/v1/analysis", handler.GetAnalysis)

	logger.Info("Service gestartet", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
