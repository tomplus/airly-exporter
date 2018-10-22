package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PromCollectors has instances of Prometheus Collectors
type PromCollectors struct {
	countTotal   prometheus.Counter
	errorTotal   prometheus.Counter
	responseTime prometheus.Histogram
	responseCode *prometheus.CounterVec
	value        *prometheus.GaugeVec
	index        *prometheus.GaugeVec
	standard     *prometheus.GaugeVec
}

// RegisterCollectors registers all collectors
func (promCollectors *PromCollectors) RegisterCollectors() {

	promCollectors.countTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "airly_count_total",
		Help: "Total number of performed check",
	})
	prometheus.MustRegister(promCollectors.countTotal)

	promCollectors.errorTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "airly_errors_total",
		Help: "Total number of errors",
	})
	prometheus.MustRegister(promCollectors.errorTotal)

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

	promCollectors.value = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_value",
		Help: "Values of the given measurement type",
	}, []string{"sensor", "name"})
	prometheus.MustRegister(promCollectors.value)

	promCollectors.index = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_index",
		Help: "Air Quality Index",
	}, []string{"sensor", "name"})
	prometheus.MustRegister(promCollectors.index)

	promCollectors.standard = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "airly_standard",
		Help: "Concentration value of a given pollutant expressed as a percentage of this concentration in the WHO standard",
	}, []string{"sensor", "name", "pollutant"})
	prometheus.MustRegister(promCollectors.standard)

}

// SetMeasurements copied values from latest measurements to Prometheus collectors
func (promCollectors *PromCollectors) SetMeasurements(sensor string, measurements MeasurementsTimeFramed) {

	for _, idx := range measurements.Values {
		sensorLabels := prometheus.Labels{"sensor": sensor, "name": idx.Name}
		promCollectors.value.With(sensorLabels).Set(idx.Value)
	}

	for _, idx := range measurements.Indexes {
		sensorLabels := prometheus.Labels{"sensor": sensor, "name": idx.Name}
		promCollectors.index.With(sensorLabels).Set(idx.Value)
	}

	for _, idx := range measurements.Standards {
		sensorLabels := prometheus.Labels{"sensor": sensor, "name": idx.Name, "pollutant": idx.Pollutant}
		promCollectors.standard.With(sensorLabels).Set(idx.Percent)
	}

}
