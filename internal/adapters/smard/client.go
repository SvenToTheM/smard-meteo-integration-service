package smard

import (
	"context"
	"time"

	"github.com/sventothem/smard-meteo-integration-service/internal/core/domain"
)

type Client struct {
	// TODO: httpClient and BaseUrl
}

func NewClient() *Client { return &Client{} }

func (c *Client) GetGridData(ctx context.Context, start, end time.Time) ([]domain.GridState, error) {
	// TODO: http request against SMARD
    return nil, nil
}
