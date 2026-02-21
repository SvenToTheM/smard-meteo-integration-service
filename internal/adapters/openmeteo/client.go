package openmeteo

import (
	"context"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

type Client struct {
	// TODO: httpClient and BaseUrl
}

func NewClient() *Client { return &Client{} }

func (c *Client) GetWeather(ctx context.Context, lat, lon float64, start, end time.Time) ([]domain.WeatherState, error) {
	// TODO: http request against Open-Meteo
    return nil, nil
}