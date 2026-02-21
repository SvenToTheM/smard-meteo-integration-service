package smard

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

func (c *Client) GetGridData(ctx context.Context, start, end time.Time) ([]domain.GridState, error) {
	// TODO: http request against SMARD
	return nil, nil
}
