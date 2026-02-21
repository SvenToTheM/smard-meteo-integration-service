package main

import (
	"log/slog"
	"net/http"
	"os"

	httpHandler "github.com/sventothem/smard-meteo-integration-service/internal/adapters/httphandler"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/openmeteo"
	"github.com/sventothem/smard-meteo-integration-service/internal/adapters/smard"
	"github.com/sventothem/smard-meteo-integration-service/internal/config"
	"github.com/sventothem/smard-meteo-integration-service/internal/core/services"
)

func main() {
	// Configure JSON structured logging for production readiness
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	extHttpClient := &http.Client{}

	sClient := smard.NewClient(cfg.SmardBaseURL)
	mClient := openmeteo.NewClient(cfg.MeteoBaseURL, extHttpClient)

	analysisService := services.NewGridService(sClient, mClient, cfg.DefaultLat, cfg.DefaultLon)

	handler := httpHandler.NewAnalysisHandler(analysisService)

	mux := http.NewServeMux()

	// register Route
	mux.HandleFunc("/api/v1/analysis", handler.GetAnalysis)

	slog.Info("Starting service", 
		slog.String("port", cfg.Port),
		slog.String("smard_url", cfg.SmardBaseURL),
		slog.String("meteo_url", cfg.MeteoBaseURL),
	)

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
