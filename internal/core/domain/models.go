package domain

import "time"

// represents the energy data of a specific point in time
type GridState struct {
	Timestamp           time.Time `json:"timestamp"`
	Consumption         float64   `json:"consumption"`
	GenerationRenewable float64   `json:"generation_renewable"`
	IsCovered           bool      `json:"is_covered"`
}

// represents the weather data of a specific point in time
type WeatherState struct {
	Timestamp      time.Time `json:"timestamp"`
	Temperature    float64   `json:"temperature"`
	WindSpeed      float64   `json:"wind_speed"`
	SolarRadiation float64   `json:"solar_radiation"`
}

// represents the enriched data of a specific point in time
type EnrichedData struct {
	GridState    GridState    `json:"grid_state"`
	WeatherState WeatherState `json:"weather_state"`
}
