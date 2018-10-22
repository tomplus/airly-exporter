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
	Current  MeasurementsTimeFramed   `json:"current,omitempty"`
	Forecast []MeasurementsTimeFramed `json:"forecast,omitempty"`
	History  []MeasurementsTimeFramed `json:"history,omitempty"`
}

// MeasurementsTimeFramed is a response from API with measurement time series
type MeasurementsTimeFramed struct {
	FromDateTime string             `json:"fromDateTime,omitempty"`
	TillDateTime string             `json:"tillDateTime,omitempty"`
	Values       []MeasuredValue    `json:"values,omitempty"`
	Indexes      []MeasuredIndex    `json:"indexes,omitempty"`
	Standards    []MeasuredStandard `json:"standards,omitempty"`
}

type MeasuredValue struct {
	Name  string  `json:"name,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type MeasuredIndex struct {
	Name        string  `json:"name,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Level       string  `json:"level,omitempty"`
	Description string  `json:"description,omitempty"`
	Advice      string  `json:"advice,omitempty"`
	Color       string  `json:"color,omitempty"`
}

type MeasuredStandard struct {
	Name      string  `json:"name,omitempty"`
	Pollutant string  `json:"pollutant,omitempty"`
	Limit     float64 `json:"limit,omitempty"`
	Percent   float64 `json:"percent,omitempty"`
}

// NewAPIClient creates a new APIClient
func NewAPIClient(apiURL string, apiKey string) *APIClient {
	return &APIClient{apiURL, apiKey}
}

// SensorMeasurements returns response from Airly API for installationId
func (api *APIClient) SensorMeasurements(installationId string) (SensorMeasurementsResponse, int, error) {

	var response SensorMeasurementsResponse

	v := url.Values{}
	v.Set("apikey", api.key)
	v.Set("installationId", installationId)
	req := api.url + "/v2/measurements/installation?" + v.Encode()

	resp, err := http.Get(req)
	if err != nil {
		return response, 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, resp.StatusCode, err
}
