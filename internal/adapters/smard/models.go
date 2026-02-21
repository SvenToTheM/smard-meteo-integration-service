package smard

type smardDataResponse struct {
	Series [][]any `json:"series"` // [[timestamp (float64), value (float64 or nil)], ...]
}
