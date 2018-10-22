# Airly-exporter

[![Build Status](https://travis-ci.org/tomplus/airly-exporter.svg?branch=master)](https://travis-ci.org/tomplus/airly-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomplus/airly-exporter)](https://goreportcard.com/report/github.com/tomplus/airly-exporter)

Airly-exporter for Prometheus.

## Overview

Airl-exporter is a server which scrapes metrics from [Airly](https://airly.eu/) and exposes them in Prometheus format. You can use
Prometheus server to scrape metrics and visualize them.

## Access to the Airly API

An API Key is required to query Airly API. You can get it for free after registration at airly.eu.
You have to know that free API token is limited to 1000 requests per day and 50 requests per minute, so
set reasonable value for the `refresh-interval` parameter.

Airly-exporter uses Airly API 2.0.

More info: [developer.airly.eu](https://developer.airly.eu/docs)

## Configuration

Airly-exporter requires parameters which can be passed via the command line, enviroment variables or a configuration file.

Available arguments as flags:

```
Usage of airly-exporter:
  -api-key string
    	Your key for Airly API
  -api-url string
    	Airly API endpoint (default "https://airapi.airly.eu")
  -config-file string
    	Path to the config file (format: flag=value\n).
  -listen-address string
    	the address to listen on for http requests. (default ":8080")
  -refresh-interval string
    	Refresh sensor interval with units (default "5m")
  -sensors string
    	Comma separated sensors IDs (default "204,822")
```

which can be replaced by enviroment variables `API_KEY`, `API_URL` etc. You can also provide configuration
via a configuration file.

```
listen-address=9090
refresh-interval=5m
sensors=204,822
```

Airl-exporter watches the configuration file and applies changes related to the `sensors` on the fly.

To get your favourite sensors IDs use Airly map, find an interesting sensor and click to see details. Sensor ID
will appear in the url (`...&id=1015`).

## Installing and running

You can use Docker to start airly-exporter:

```
docker pull tpimages/airly-exporter:latest
docker run --rm -e API_KEY=my-api-key -p 8080:8080 tpimages/airly-exporter:latest
```

and metrics are exposed via http://localhost:8080/metrics

Alternatively you can install this using `go`:

```go get github.com/tomplus/airly-exporter```

or download binary file from [airly-exporter/releases](https://github.com/tomplus/airly-exporter/releases).

## Running on Kubernetes cluster with Prometheus Operator

The repository contains example manifests to deploy Airly-exporter to Kubernetes with
[Prometheus Operator](https://github.com/coreos/prometheus-operator) installed. There are manifests
for creating Deployment, Service and Service Monitor.

## List of exposed metrics:

```
# HELP airly_count_total Total number of performed check
# TYPE airly_count_total counter
airly_count_total 2
# HELP airly_errors_total Total number of errors
# TYPE airly_errors_total counter
airly_errors_total 0
# HELP airly_index Air Quality Index
# TYPE airly_index gauge
airly_index{name="AIRLY_CAQI",sensor="204"} 76.68
airly_index{name="AIRLY_CAQI",sensor="822"} 27
# HELP airly_request_duration_seconds Histogram of request duration
# TYPE airly_request_duration_seconds histogram
airly_request_duration_seconds_bucket{le="0.005"} 0
airly_request_duration_seconds_bucket{le="0.01"} 0
airly_request_duration_seconds_bucket{le="0.025"} 0
airly_request_duration_seconds_bucket{le="0.05"} 0
airly_request_duration_seconds_bucket{le="0.1"} 0
airly_request_duration_seconds_bucket{le="0.25"} 0
airly_request_duration_seconds_bucket{le="0.5"} 0
airly_request_duration_seconds_bucket{le="1"} 2
airly_request_duration_seconds_bucket{le="2.5"} 2
airly_request_duration_seconds_bucket{le="5"} 2
airly_request_duration_seconds_bucket{le="10"} 2
airly_request_duration_seconds_bucket{le="+Inf"} 2
airly_request_duration_seconds_sum 1.381165496
airly_request_duration_seconds_count 2
# HELP airly_response_code Response code from Airly API
# TYPE airly_response_code counter
airly_response_code{code="200"} 2
# HELP airly_standard Concentration value of a given pollutant expressed as a percentage of this concentration in the WHO standard
# TYPE airly_standard gauge
airly_standard{name="WHO",pollutant="PM10",sensor="204"} 192.07
airly_standard{name="WHO",pollutant="PM10",sensor="822"} 46.68
airly_standard{name="WHO",pollutant="PM25",sensor="204"} 230.71
airly_standard{name="WHO",pollutant="PM25",sensor="822"} 64.79
# HELP airly_value Values of the given measurement type
# TYPE airly_value gauge
airly_value{name="HUMIDITY",sensor="204"} 87.69
airly_value{name="HUMIDITY",sensor="822"} 94.57
airly_value{name="PM1",sensor="204"} 37.21
airly_value{name="PM1",sensor="822"} 9.98
airly_value{name="PM10",sensor="204"} 96.04
airly_value{name="PM10",sensor="822"} 23.34
airly_value{name="PM25",sensor="204"} 57.68
airly_value{name="PM25",sensor="822"} 16.2
airly_value{name="PRESSURE",sensor="204"} 1024.41
airly_value{name="PRESSURE",sensor="822"} 1028.68
airly_value{name="TEMPERATURE",sensor="204"} 9.23
airly_value{name="TEMPERATURE",sensor="822"} 5.81
```

... plus a bunch of metrics from the prometheus client.

