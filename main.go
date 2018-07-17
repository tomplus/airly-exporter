// Copyright 2018 Airly-exporter Authors

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AirlyExporter struct {
	configFile      *string
	listenAddress   *string
	apiUrl          *string
	apiKey          *string
	refreshInterval *string
	sensors         *string
	promCollectors  PromCollectors
	api             *ApiClient
}

func FlagStringWithDefaultFromEnv(name string, value string, usage string) *string {
	env_name := strings.Replace(strings.ToUpper(name), "-", "_", -1)
	def_value := os.Getenv(env_name)
	if def_value == "" {
		def_value = value
	}
	return flag.String(name, def_value, usage)
}

func (airlyExporter AirlyExporter) ReadConfigFile() {
	log.Printf("Load configuration from file %v", *airlyExporter.configFile)
	file, err := os.Open(*airlyExporter.configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.SplitN(line, "=", 2)
		f := flag.Lookup(kv[0])
		if f == nil {
			log.Printf("Ignore line %v", line)
		} else {
			flag.Set(kv[0], kv[1])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (airlyExporter AirlyExporter) WatchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(*airlyExporter.configFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				airlyExporter.ReadConfigFile()
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (airlyExporter AirlyExporter) WatchSensors() {

	duration, err := time.ParseDuration(*airlyExporter.refreshInterval)
	if err != nil {
		log.Fatal(err)
	}
	airlyExporter.api = NewApiClient(*airlyExporter.apiUrl, *airlyExporter.apiKey)

	for {
		sensors := strings.Split(*airlyExporter.sensors, ",")
		for _, sensor := range sensors {
			log.Printf("Get data from %v\n", sensor)
			airlyExporter.QuerySensor(sensor)
		}
		<-time.After(duration)
	}
}

func (airlyExporter AirlyExporter) QuerySensor(sensor string) {

	airlyExporter.promCollectors.count_total.Inc()

	timer := prometheus.NewTimer(airlyExporter.promCollectors.responseTime)
	defer timer.ObserveDuration()

	measurements, code, err := airlyExporter.api.SensorMeasurements(sensor)
	if code > 0 {
		airlyExporter.promCollectors.responseCode.With(prometheus.Labels{"code": strconv.Itoa(code)}).Inc()
	}
	if err != nil {
		log.Println(err)
		airlyExporter.promCollectors.error_total.Inc()
		return
	}
	if code != 200 {
		fmt.Println("Unexpected response code %v", code)
		return
	}

	fmt.Printf("Debug measurements %+v\n", measurements.CurrentMeasurements)
	airlyExporter.promCollectors.SetMeasurements(sensor, measurements.CurrentMeasurements)
}

func (airlyExporter *AirlyExporter) FlagParse() {

	airlyExporter.configFile = FlagStringWithDefaultFromEnv("config-file", "", "Path to the config file (format: flag=value\\n).")
	airlyExporter.listenAddress = FlagStringWithDefaultFromEnv("listen-address", ":8080", "the address to listen on for http requests.")
	airlyExporter.apiKey = FlagStringWithDefaultFromEnv("api-key", "", "Your key for Airly API")
	airlyExporter.apiUrl = FlagStringWithDefaultFromEnv("api-url", "https://airapi.airly.eu", "Airly API endpoint")
	airlyExporter.refreshInterval = FlagStringWithDefaultFromEnv("refresh-interval", "5m", "Refresh sensor interval with units")
	airlyExporter.sensors = FlagStringWithDefaultFromEnv("sensors", "204,822", "Comma separated sensors IDs")
	flag.Parse()
}

func main() {

	airlyExporter := AirlyExporter{}

	airlyExporter.FlagParse()
	airlyExporter.promCollectors.RegisterCollectors()

	if *airlyExporter.configFile != "" {
		airlyExporter.ReadConfigFile()
		go airlyExporter.WatchConfig()
	}

	log.Println("Airly-exporter started")

	go airlyExporter.WatchSensors()

	http.Handle("/", http.RedirectHandler("/metrics", 302))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*airlyExporter.listenAddress, nil))
}
