package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// APIClient manages communication with Arily API
type APIClient struct {
	url string
	key string
}

// SensorMeasurementsResponse is a response from API with measurements
type SensorMeasurementsResponse struct {
	CurrentMeasurements AllMeasurements          `json:"currentMeasurements,omitempty"`
	Forecast            []MeasurementsTimeFramed `json:"forecast,omitempty"`
	History             []MeasurementsTimeFramed `json:"history,omitempty"`
}

// MeasurementsTimeFramed is a response from API with measurement time series
type MeasurementsTimeFramed struct {
	FromDateTime string          `json:"fromDateTime,omitempty"`
	Measurements AllMeasurements `json:"measurements,omitempty"`
	TillDateTime string          `json:"tillDateTime,omitempty"`
}

// AllMeasurements is a response from API with one set of measurement
type AllMeasurements struct {
	AirQualityIndex float64 `json:"airQualityIndex,omitempty"`
	Humidity        float64 `json:"humidity,omitempty"`
	MeasurementTime string  `json:"measurementTime,omitempty"`
	Pm1             float64 `json:"pm1,omitempty"`
	Pm10            float64 `json:"pm10,omitempty"`
	Pm25            float64 `json:"pm25,omitempty"`
	PollutionLevel  float64 `json:"pollutionLevel,omitempty"`
	Pressure        float64 `json:"pressure,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
	WindDirection   float64 `json:"windDirection,omitempty"`
	WindSpeed       float64 `json:"windSpeed,omitempty"`
}

// NewAPIClient creates a new APIClient
func NewAPIClient(apiURL string, apiKey string) *APIClient {
	return &APIClient{apiURL, apiKey}
}

// SensorMeasurements returns response from Airly API for sensorID
func (api *APIClient) SensorMeasurements(sensorID string) (SensorMeasurementsResponse, int, error) {

	var response SensorMeasurementsResponse

	v := url.Values{}
	v.Set("apikey", api.key)
	v.Set("sensorId", sensorID)
	req := api.url + "/v1/sensor/measurements?" + v.Encode()

	resp, err := http.Get(req)
	if err != nil {
		return response, 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, resp.StatusCode, err
}
