package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const TestResponseValid = `{
  "currentMeasurements": {
    "airQualityIndex": 1,
    "humidity": 2,
    "measurementTime": "string",
    "pm1": 3,
    "pm10": 4,
    "pm25": 5,
    "pollutionLevel": 6,
    "pressure": 7,
    "temperature": 8,
    "windDirection": 9,
    "windSpeed": 10
  },
  "forecast": [
    {
      "fromDateTime": "string",
      "measurements": {
        "airQualityIndex": 11,
        "humidity": 12,
        "measurementTime": "string",
        "pm1": 13,
        "pm10": 14,
        "pm25": 15,
        "pollutionLevel": 16,
        "pressure": 17,
        "temperature": 18,
        "windDirection": 19,
        "windSpeed": 20
      },
      "tillDateTime": "string"
    }
  ],
  "history": [
    {
      "fromDateTime": "string",
      "measurements": {
        "airQualityIndex": 21,
        "humidity": 22,
        "measurementTime": "string",
        "pm1": 23,
        "pm10": 24,
        "pm25": 25,
        "pollutionLevel": 26,
        "pressure": 27,
        "temperature": 28,
        "windDirection": 29,
        "windSpeed": 30
      },
      "tillDateTime": "string"
    }
  ]
}`

func TestSensorMeasurements(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, TestResponseValid)
	}))
	defer ts.Close()

	api := NewAPIClient(ts.URL, "API-KEY")
	resp, code, err := api.SensorMeasurements("1234")

	if err != nil {
		t.Errorf("assert error %v is not nil", err)
	}

	if code != 200 {
		t.Errorf("assert http code %v is not 200", code)
	}

	expected := SensorMeasurementsResponse{
		CurrentMeasurements: AllMeasurements{AirQualityIndex: 1,
			Humidity:        2,
			MeasurementTime: "string",
			Pm1:             3,
			Pm10:            4,
			Pm25:            5,
			PollutionLevel:  6,
			Pressure:        7,
			Temperature:     8,
			WindDirection:   9,
			WindSpeed:       10},
	}

	if !reflect.DeepEqual(resp.CurrentMeasurements, expected.CurrentMeasurements) {
		t.Errorf("assert response %+v is not %v", resp, expected)
	}
}
