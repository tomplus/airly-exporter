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

// AirlyExporter contains parameters for Airly-exporter
type AirlyExporter struct {
	configFile      *string
	listenAddress   *string
	apiURL          *string
	apiKey          *string
	refreshInterval *string
	sensors         *string
	promCollectors  PromCollectors
	api             *APIClient
}

func flagStringWithDefaultFromEnv(name string, value string, usage string) *string {
	envName := strings.Replace(strings.ToUpper(name), "-", "_", -1)
	defValue := os.Getenv(envName)
	if defValue == "" {
		defValue = value
	}
	return flag.String(name, defValue, usage)
}

func (airlyExporter AirlyExporter) readConfigFile() {
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

func (airlyExporter AirlyExporter) watchConfig() {
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
				airlyExporter.readConfigFile()
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (airlyExporter AirlyExporter) watchSensors() {

	duration, err := time.ParseDuration(*airlyExporter.refreshInterval)
	if err != nil {
		log.Fatal(err)
	}
	airlyExporter.api = NewAPIClient(*airlyExporter.apiURL, *airlyExporter.apiKey)

	for {
		sensors := strings.Split(*airlyExporter.sensors, ",")
		for _, sensor := range sensors {
			log.Printf("Get data from %v\n", sensor)
			airlyExporter.querySensor(sensor)
		}
		<-time.After(duration)
	}
}

func (airlyExporter AirlyExporter) querySensor(sensor string) {

	airlyExporter.promCollectors.countTotal.Inc()

	timer := prometheus.NewTimer(airlyExporter.promCollectors.responseTime)
	defer timer.ObserveDuration()

	measurements, code, err := airlyExporter.api.SensorMeasurements(sensor)
	if code > 0 {
		airlyExporter.promCollectors.responseCode.With(prometheus.Labels{"code": strconv.Itoa(code)}).Inc()
	}
	if err != nil {
		log.Println(err)
		airlyExporter.promCollectors.errorTotal.Inc()
		return
	}
	if code != 200 {
		log.Printf("Unexpected response code %v", code)
		return
	}

	fmt.Printf("Debug measurements %+v\n", measurements.Current)
	airlyExporter.promCollectors.SetMeasurements(sensor, measurements.Current)
}

func (airlyExporter *AirlyExporter) flagParse() {

	airlyExporter.configFile = flagStringWithDefaultFromEnv("config-file", "", "Path to the config file (format: flag=value\\n).")
	airlyExporter.listenAddress = flagStringWithDefaultFromEnv("listen-address", ":8080", "the address to listen on for http requests.")
	airlyExporter.apiKey = flagStringWithDefaultFromEnv("api-key", "", "Your key for Airly API")
	airlyExporter.apiURL = flagStringWithDefaultFromEnv("api-url", "https://airapi.airly.eu", "Airly API endpoint")
	airlyExporter.refreshInterval = flagStringWithDefaultFromEnv("refresh-interval", "5m", "Refresh sensor interval with units")
	airlyExporter.sensors = flagStringWithDefaultFromEnv("sensors", "204,822", "Comma separated sensors/installations IDs")
	flag.Parse()
}

func main() {

	airlyExporter := AirlyExporter{}

	airlyExporter.flagParse()
	airlyExporter.promCollectors.RegisterCollectors()

	if *airlyExporter.configFile != "" {
		airlyExporter.readConfigFile()
		go airlyExporter.watchConfig()
	}

	log.Println("Airly-exporter started")

	go airlyExporter.watchSensors()

	http.Handle("/", http.RedirectHandler("/metrics", 302))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*airlyExporter.listenAddress, nil))
}
