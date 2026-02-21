package services

import (
	"context"
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
	// get energy data
	energyData, err := s.energyProvider.GetGridData(ctx, start, end)
	if err != nil {
		//TODO: log if energy is null or err
		return nil, err
	}

	// get weather data
	weatherData, err := s.weatherProvider.GetWeather(ctx, s.lat, s.lon, start, end)
	// TODO: log if weather is null

	return s.mergeAndCheck(energyData, weatherData), nil
}

func (s *GridService) mergeAndCheck(energy []domain.GridState, weather []domain.WeatherState) []domain.EnrichedData {

	// TODO:Combine  & check if Gen >= Cons
	return nil
}
