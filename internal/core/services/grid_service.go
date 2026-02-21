package services

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
	"github.com/sventothem/smard-meteo-integration-service/internal/core/ports"
)

type GridService struct {
	energyProvider  ports.EnergyProvider
	weatherProvider ports.WeatherProvider
	lat             float64
	lon             float64
}

func NewGridService(e ports.EnergyProvider, w ports.WeatherProvider, lat, lon float64) *GridService {
	return &GridService{
		energyProvider:  e,
		weatherProvider: w,
		lat:             lat,
		lon:             lon,
	}
}

func (s *GridService) Analyze(ctx context.Context, start, end time.Time) ([]domain.EnrichedData, error) {
	slog.Info("Fetching energy data", slog.Time("start", start), slog.Time("end", end))
	energyData, err := s.energyProvider.GetGridData(ctx, start, end)
	if err != nil {
		slog.Error("Failed to fetch energy data", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Fetching weather data", slog.Float64("lat", s.lat), slog.Float64("lon", s.lon))
	weatherData, err := s.weatherProvider.GetWeatherData(ctx, s.lat, s.lon, start, end)
	if err != nil {
		slog.Warn("Failed to fetch weather data, continuing without weather context", slog.Any("error", err))
		// Resilience: continuing without weather data
	}

	slog.Info("Merging data and checking coverage", slog.Int("energy_points", len(energyData)), slog.Int("weather_points", len(weatherData)))
	return s.mergeAndCheck(energyData, weatherData), nil
}

func (s *GridService) mergeAndCheck(energy []domain.GridState, weather []domain.WeatherState) []domain.EnrichedData {
	weatherMap := make(map[int64]*domain.WeatherState, len(weather))
	// use timestamp as key
	for i := range weather {
		w := &weather[i]
		weatherMap[w.Timestamp.Unix()] = w
	}

	var results []domain.EnrichedData

	for _, e := range energy {
		// Core logic, if generation >= consumption
		// TODO: Optional http path param
		if e.GenerationRenewable >= e.Consumption {

			enriched := domain.EnrichedData{
				GridState: e,
			}

			// Add weather data if timestamp is avaiable
			if w, ok := weatherMap[e.Timestamp.Unix()]; ok {
				enriched.WeatherState = w
			}

			results = append(results, enriched)
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].GridState.Timestamp.After(results[j].GridState.Timestamp)
		})
	}

	return results
}
