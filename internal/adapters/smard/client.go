package smard

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

const (
	FilterConsumption   = "410"
	FilterWindOffshore  = "1225"
	FilterWindOnshore   = "4067"
	FilterSolar         = "4068"
	FilterBiomass       = "4066"
	FilterHydro         = "1226"
	FilterOtherRenew    = "1228"
)

var renewableFilters = []string{
	FilterWindOffshore, FilterWindOnshore, FilterSolar,
	FilterBiomass, FilterHydro, FilterOtherRenew,
}

type smardIndexResponse struct {
	Timestamps []int64 `json:"timestamps"`
}

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

// fetches and aggregates consumption and renewable generation
func (c *Client) GetGridData(ctx context.Context, start, end time.Time) ([]domain.GridState, error) {
	
	// Resolve correct file timestamp from SMARD Index
	fileTS, err := c.getRelevantTimestamp(ctx, FilterConsumption, start)

	if err != nil {
		slog.Error("Failed to resolve SMARD timestamp index", slog.Any("error", err))
		return nil, fmt.Errorf("failed to fetch timestamp index: %w", err)
	}

	slog.Debug("Starting to fetch SMARD load profile", slog.String("filter", FilterConsumption))
	
	// Fetch load profile: Consumption
	consumption, err := c.fetch(ctx, FilterConsumption, fileTS)

	if err != nil {
		slog.Error("Critical failure: unable to fetch consumption data", slog.Any("error", err))
		return nil, fmt.Errorf("failed to fetch consumption: %w", err)
	}

	slog.Debug("Fetching renewable generation data", slog.Int("filter_count", len(renewableFilters)))
	
	// Fetch and aggregate all renewable sources
	totalRenewables := make(map[int64]float64)

	for _, filter := range renewableFilters {
		data, err := c.fetch(ctx, filter, fileTS)
		if err != nil {
			slog.Warn("Skipping renewable source due to fetch error", slog.String("filter", filter), slog.Any("error", err))
			// Resilience: If one source fails (e.g., biomass API is down), 
			// we skip it but keep the rest of the generation data.
			continue 
		}
		
		// Sum up values for each timestamp
		for ts, val := range data {
			totalRenewables[ts] += val
		}
	}

	var results []domain.GridState

	// Merge consumption and aggregated generation
	for ts, consVal := range consumption {
		t := time.UnixMilli(ts).UTC()
		
		// Filter data points outside the requested time window
		if t.Before(start) || t.After(end) {
			continue
		}

		gen := totalRenewables[ts] 
		
		results = append(results, domain.GridState{
			Timestamp:           t,
			Consumption:         consVal,
			GenerationRenewable: gen,
			IsCovered:           gen >= consVal, // Core business logic
		})
	}
	slog.Info("results", results)
	return results, nil
}

// getRelevantTimestamp fetches index_hour.json to find the 
// start timestamp of the file containing the requested date
func (c *Client) getRelevantTimestamp(ctx context.Context, filter string, target time.Time) (int64, error) {
	url := fmt.Sprintf("%s/chart_data/%s/DE/index_hour.json", c.baseURL, filter)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("index api returned status %d", resp.StatusCode)
	}

	var index smardIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return 0, err
	}

	targetMilli := target.UnixMilli()

	var closestTS int64

	for _, ts := range index.Timestamps {
		if ts <= targetMilli {
			closestTS = ts
		} else {
			break
		}
	}

	if closestTS == 0 && len(index.Timestamps) > 0 {
		return index.Timestamps[0], nil
	}
	return closestTS, nil
}

// fetch performs the actual HTTP request for a specific filter and timestamp
func (c *Client) fetch(ctx context.Context, filter string, fileTS int64) (map[int64]float64, error) {
	url := fmt.Sprintf("%s/chart_data/%s/DE/%s_DE_hour_%d.json", c.baseURL, filter, filter, fileTS)

	slog.Debug("Executing SMARD HTTP request", slog.String("url", url), slog.String("filter", filter))
	
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := c.httpClient.Do(req)
	
	if err != nil {
		slog.Error("HTTP request failed", slog.String("url", url), slog.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var raw smardDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	data := make(map[int64]float64)
	for _, p := range raw.Series {
		// SMARD format [[timestamp, value], ...]
		if len(p) == 2 && p[1] != nil {
			ts := int64(p[0].(float64))
			val := p[1].(float64)
			data[ts] = val
		}
	}
	return data, nil
}