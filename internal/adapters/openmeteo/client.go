package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) GetWeatherData(ctx context.Context, lat, lon float64, start, end time.Time) ([]domain.WeatherState, error) {
	layout := "2006-01-02"

	startDateStr := start.Format(layout)
	endDateStr := end.Format(layout)

	url := fmt.Sprintf(
		"%s"+
			"?latitude=%.4f"+
			"&longitude=%.4f"+
			"&start_date=%s"+
			"&end_date=%s"+
			"&hourly=temperature_2m,wind_speed_100m,shortwave_radiation",
		c.baseURL, lat, lon, startDateStr, endDateStr)

	slog.Debug("Executing Open-Meteo Http reqeuest", slog.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create weather request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("HTTP request to Open-Meteo failed", slog.String("url", url), slog.Any("error", err))
		return nil, fmt.Errorf("weather api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Open-Meteo API returned non-OK status", slog.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("weather api returned status %d", resp.StatusCode)
	}

	var raw openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("Failed to decode Open-Meteo response", slog.Any("error", err))
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	// Pre-allocate slice for performance
	capacity := len(raw.Hourly.Time)
	results := make([]domain.WeatherState, 0, capacity)

	// Map raw JSON arrays to domain objects
	for i, tStr := range raw.Hourly.Time {
		// Open-Meteo timestamp format is "2006-01-02T15:04"
		t, err := time.Parse("2006-01-02T15:04", tStr)
		if err != nil {
			slog.Warn("Skipping invalid timestamp from weather API", slog.String("timestamp", tStr))
			continue
		}

		results = append(results, domain.WeatherState{
			Timestamp:      t.UTC(),
			Temperature:    raw.Hourly.Temperature2m[i],
			WindSpeed:      raw.Hourly.WindSpeed100m[i],
			SolarRadiation: raw.Hourly.ShortwaveRadiation[i],
		})
	}

	slog.Debug("Successfully fetched weather data", slog.Int("count", len(results)))
	return results, nil
}
