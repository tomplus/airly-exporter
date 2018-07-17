package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PromCollectors struct {
	count_total     prometheus.Counter
	error_total     prometheus.Counter
	responseTime    prometheus.Histogram
	responseCode    *prometheus.CounterVec
	airQualityIndex *prometheus.GaugeVec
	humidity        *prometheus.GaugeVec
	pm1             *prometheus.GaugeVec
	pm10            *prometheus.GaugeVec
	pm25            *prometheus.GaugeVec
	pollutionLevel  *prometheus.GaugeVec
	pressure        *prometheus.GaugeVec
	temperature     *prometheus.GaugeVec
	windDirection   *prometheus.GaugeVec
	windSpeed       *prometheus.GaugeVec
}

func (promCollectors *PromCollectors) RegisterCollectors() {

	promCollectors.count_total = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "airly_count_total",
		Help: "Total number of performed check",
	})
	prometheus.MustRegister(promCollectors.count_total)

	promCollectors.error_total = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "airly_errors_total",
		Help: "Total number of errors",
	})
	prometheus.MustRegister(promCollectors.error_total)

	promCollectors.responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "airly_request_duration_seconds",
		Help: "Histogram of request duration",
	})

	prometheus.MustRegister(promCollectors.responseTime)

	promCollectors.responseCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "airly_response_code",
		Help: "Response code from Airly API",
	}, []string{"code"})
	prometheus.MustRegister(promCollectors.responseCode)

	promCollectors.airQualityIndex = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_air_quality_index",
		Help: "Common Air Quality Index (CAQI)",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.airQualityIndex)

	promCollectors.humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_humidity",
		Help: "Humidity",
	}, []string{"sensor"})

	prometheus.MustRegister(promCollectors.humidity)

	promCollectors.pm1 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_pm_1",
		Help: "PM1",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.pm1)

	promCollectors.pm10 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_pm_10",
		Help: "PM10",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.pm10)

	promCollectors.pm25 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_pm_25",
		Help: "PM25",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.pm25)

	promCollectors.pollutionLevel = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_pollution_level",
		Help: "Pollution level based on CAQI value. Possible values: [0 to 6]. 0 - unknown, 1 - best air, 6 - worst",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.pollutionLevel)

	promCollectors.pressure = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_pressure",
		Help: "Pressure",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.pressure)

	promCollectors.temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_temperature",
		Help: "Temperature",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.temperature)

	promCollectors.windDirection = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_wind_direction",
		Help: "Wind direction",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.windDirection)

	promCollectors.windSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_wind_speed",
		Help: "Wind speed",
	}, []string{"sensor"})
	prometheus.MustRegister(promCollectors.windSpeed)
}

func (promCollectors *PromCollectors) SetMeasurements(sensor string, measurements AllMeasurements) {

	sensor_labels := prometheus.Labels{"sensor": sensor}

	promCollectors.airQualityIndex.With(sensor_labels).Set(measurements.AirQualityIndex)
	promCollectors.humidity.With(sensor_labels).Set(measurements.Humidity)
	promCollectors.pm1.With(sensor_labels).Set(measurements.Pm1)
	promCollectors.pm10.With(sensor_labels).Set(measurements.Pm10)
	promCollectors.pm25.With(sensor_labels).Set(measurements.Pm25)
	promCollectors.pollutionLevel.With(sensor_labels).Set(measurements.PollutionLevel)
	promCollectors.pressure.With(sensor_labels).Set(measurements.Pressure)
	promCollectors.temperature.With(sensor_labels).Set(measurements.Temperature)
	promCollectors.windDirection.With(sensor_labels).Set(measurements.WindDirection)
	promCollectors.windSpeed.With(sensor_labels).Set(measurements.WindSpeed)
}
