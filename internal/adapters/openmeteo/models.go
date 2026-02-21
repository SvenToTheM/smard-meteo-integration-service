package openmeteo

type openMeteoResponse struct {
	Hourly struct {
		Time               []string  `json:"time"`
		Temperature2m      []float64 `json:"temperature_2m"`
		WindSpeed100m      []float64 `json:"wind_speed_100m"`
		ShortwaveRadiation []float64 `json:"shortwave_radiation"`
	} `json:"hourly"`
}
