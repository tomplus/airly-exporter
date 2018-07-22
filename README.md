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
# HELP airly_air_quality_index Common Air Quality Index (CAQI)
# TYPE airly_air_quality_index gauge
airly_air_quality_index{sensor="204"} 33.64089591836734
airly_air_quality_index{sensor="822"} 21.70154666666667
# HELP airly_count_total Total number of performed check
# TYPE airly_count_total counter
airly_count_total 2
# HELP airly_errors_total Total number of errors
# TYPE airly_errors_total counter
airly_errors_total 0
# HELP airly_humidity Humidity
# TYPE airly_humidity gauge
airly_humidity{sensor="204"} 94.75
airly_humidity{sensor="822"} 81.40867607136424
# HELP airly_pm_1 PM1
# TYPE airly_pm_1 gauge
airly_pm_1{sensor="204"} 17.1823387755102
airly_pm_1{sensor="822"} 9.577728
# HELP airly_pm_10 PM10
# TYPE airly_pm_10 gauge
airly_pm_10{sensor="204"} 38.978979591836726
airly_pm_10{sensor="822"} 17.10304
# HELP airly_pm_25 PM25
# TYPE airly_pm_25 gauge
airly_pm_25{sensor="204"} 20.184537551020405
airly_pm_25{sensor="822"} 13.020928000000001
# HELP airly_pollution_level Pollution level based on CAQI value. Possible values: [0 to 6]. 0 - unknown, 1 - best air, 6 - worst
# TYPE airly_pollution_level gauge
airly_pollution_level{sensor="204"} 2
airly_pollution_level{sensor="822"} 1
# HELP airly_pressure Pressure
# TYPE airly_pressure gauge
airly_pressure{sensor="204"} 100935.324470464
airly_pressure{sensor="822"} 101197.11217548826
# HELP airly_request_duration_seconds Histogram of request duration
# TYPE airly_request_duration_seconds histogram
airly_request_duration_seconds_bucket{le="0.005"} 0
airly_request_duration_seconds_bucket{le="0.01"} 0
airly_request_duration_seconds_bucket{le="0.025"} 0
airly_request_duration_seconds_bucket{le="0.05"} 0
airly_request_duration_seconds_bucket{le="0.1"} 0
airly_request_duration_seconds_bucket{le="0.25"} 1
airly_request_duration_seconds_bucket{le="0.5"} 2
airly_request_duration_seconds_bucket{le="1"} 2
airly_request_duration_seconds_bucket{le="2.5"} 2
airly_request_duration_seconds_bucket{le="5"} 2
airly_request_duration_seconds_bucket{le="10"} 2
airly_request_duration_seconds_bucket{le="+Inf"} 2
airly_request_duration_seconds_sum 0.506343012
airly_request_duration_seconds_count 2
# HELP airly_response_code Response code from Airly API
# TYPE airly_response_code counter
airly_response_code{code="200"} 2
# HELP airly_temperature Temperature
# TYPE airly_temperature gauge
airly_temperature{sensor="204"} 18.073265306122448
airly_temperature{sensor="822"} 18.1751
# HELP airly_wind_direction Wind direction
# TYPE airly_wind_direction gauge
airly_wind_direction{sensor="204"} 0
airly_wind_direction{sensor="822"} 0
# HELP airly_wind_speed Wind speed
# TYPE airly_wind_speed gauge
airly_wind_speed{sensor="204"} 0
airly_wind_speed{sensor="822"} 0
```

... plus a bunch of metrics from prometheus clients.

