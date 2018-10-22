package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const TestResponseValid = `{
  "current": {
    "fromDateTime": "2018-01-01T21:46:25.476Z",
    "tillDateTime": "2018-01-01T22:46:25.476Z",
    "values": [
      {
        "name": "PM1",
        "value": 17.01
      },
      {
        "name": "PM25",
        "value": 25.69
      },
      {
        "name": "PM10",
        "value": 49.7
      },
      {
        "name": "PRESSURE",
        "value": 1025.82
      },
      {
        "name": "HUMIDITY",
        "value": 82.91
      },
      {
        "name": "TEMPERATURE",
        "value": 7.64
      }
    ],
    "indexes": [
      {
        "name": "AIRLY_CAQI",
        "value": 49.7,
        "level": "LOW",
        "description": "Air is quite good.",
        "advice": "How about going for a walk?",
        "color": "#D1CF1E"
      }
    ],
    "standards": [
      {
        "name": "WHO",
        "pollutant": "PM25",
        "limit": 25,
        "percent": 102.77
      }
    ]
  },
  "history": [
    {
      "fromDateTime": "2018-10-19T22:00:00Z",
      "tillDateTime": "2018-10-19T23:00:00Z",
      "values": [
        {
          "name": "PM1",
          "value": 26.66
        },
        {
          "name": "PM25",
          "value": 42.97
        },
        {
          "name": "PM10",
          "value": 79.19
        },
        {
          "name": "PRESSURE",
          "value": 1025.44
        },
        {
          "name": "HUMIDITY",
          "value": 92.64
        },
        {
          "name": "TEMPERATURE",
          "value": 10.22
        }
      ],
      "indexes": [
        {
          "name": "AIRLY_CAQI",
          "value": 68.24,
          "level": "MEDIUM",
          "description": "Well... It's been better.",
          "advice": "Neither good nor bad. Think before leaving the house.",
          "color": "#EFBB0F"
        }
      ],
      "standards": [
        {
          "name": "WHO",
          "pollutant": "PM25",
          "limit": 25,
          "percent": 171.88
        }
      ]
    }
  ],
  "forecast": [
    {
      "fromDateTime": "2018-01-01T22:00:00Z",
      "tillDateTime": "2018-01-01T23:00:00Z",
      "values": [
        {
          "name": "PM25",
          "value": 27.48
        },
        {
          "name": "PM10",
          "value": 53.29
        }
      ],
      "indexes": [
        {
          "name": "AIRLY_CAQI",
          "value": 52.06,
          "level": "MEDIUM",
          "description": "Well... It's been better.",
          "advice": "Protect your lungs!",
          "color": "#EFBB0F"
        }
      ],
      "standards": [
        {
          "name": "WHO",
          "pollutant": "PM25",
          "limit": 25,
          "percent": 109.93
        }
      ]
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
		Current: MeasurementsTimeFramed{
			FromDateTime: "2018-01-01T21:46:25.476Z",
			TillDateTime: "2018-01-01T22:46:25.476Z",
			Values: []MeasuredValue{
				MeasuredValue{Name: "PM1", Value: 17.01},
				MeasuredValue{Name: "PM25", Value: 25.69},
				MeasuredValue{Name: "PM10", Value: 49.7},
				MeasuredValue{Name: "PRESSURE", Value: 1025.82},
				MeasuredValue{Name: "HUMIDITY", Value: 82.91},
				MeasuredValue{Name: "TEMPERATURE", Value: 7.64}},
			Indexes: []MeasuredIndex{
				MeasuredIndex{
					Name:        "AIRLY_CAQI",
					Value:       49.7,
					Level:       "LOW",
					Description: "Air is quite good.",
					Advice:      "How about going for a walk?",
					Color:       "#D1CF1E"}},
			Standards: []MeasuredStandard{
				MeasuredStandard{
					Name:      "WHO",
					Pollutant: "PM25",
					Limit:     25,
					Percent:   102.77}},
		},
		History: []MeasurementsTimeFramed{
			MeasurementsTimeFramed{
				FromDateTime: "2018-10-19T22:00:00Z",
				TillDateTime: "2018-10-19T23:00:00Z",
				Values: []MeasuredValue{
					MeasuredValue{Name: "PM1", Value: 26.66},
					MeasuredValue{Name: "PM25", Value: 42.97},
					MeasuredValue{Name: "PM10", Value: 79.19},
					MeasuredValue{Name: "PRESSURE", Value: 1025.44},
					MeasuredValue{Name: "HUMIDITY", Value: 92.64},
					MeasuredValue{Name: "TEMPERATURE", Value: 10.22}},
				Indexes: []MeasuredIndex{
					MeasuredIndex{
						Name:        "AIRLY_CAQI",
						Value:       68.24,
						Level:       "MEDIUM",
						Description: "Well... It's been better.",
						Advice:      "Neither good nor bad. Think before leaving the house.",
						Color:       "#EFBB0F"}},
				Standards: []MeasuredStandard{
					MeasuredStandard{
						Name:      "WHO",
						Pollutant: "PM25",
						Limit:     25,
						Percent:   171.88}},
			},
		},
		Forecast: []MeasurementsTimeFramed{
			MeasurementsTimeFramed{
				FromDateTime: "2018-01-01T22:00:00Z",
				TillDateTime: "2018-01-01T23:00:00Z",
				Values: []MeasuredValue{
					MeasuredValue{Name: "PM25", Value: 27.48},
					MeasuredValue{Name: "PM10", Value: 53.29}},
				Indexes: []MeasuredIndex{
					MeasuredIndex{
						Name:        "AIRLY_CAQI",
						Value:       52.06,
						Level:       "MEDIUM",
						Description: "Well... It's been better.",
						Advice:      "Protect your lungs!",
						Color:       "#EFBB0F"}},
				Standards: []MeasuredStandard{
					MeasuredStandard{
						Name:      "WHO",
						Pollutant: "PM25",
						Limit:     25,
						Percent:   109.93}},
			}},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Errorf("assert response %+v is not %+v", resp, expected)
	}
}
