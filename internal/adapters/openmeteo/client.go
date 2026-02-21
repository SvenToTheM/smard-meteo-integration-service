package openmeteo

import (
	"context"
	"net/http"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		httpClient: client,
		baseURL:    baseURL,
	}
}

func (c *Client) GetWeather(ctx context.Context, lat, lon float64, start, end time.Time) ([]domain.WeatherState, error) {
	// TODO: http request against Open-Meteo
	return nil, nil
}