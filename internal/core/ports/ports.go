package ports

import (
	"context"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

type EnergyProvider interface {
	GetGridData(ctx context.Context, start, end time.Time) ([]domain.GridState, error)
}

type WeatherProvider interface {
	GetWeatherData(ctx context.Context, lat, lon float64, start, end time.Time) ([]domain.WeatherState, error)
}

type AnalysisService interface {
	Analyze(ctx context.Context, start, end time.Time) ([]domain.EnrichedData, error)
}
