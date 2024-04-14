package main

import (
	"crypto/tls"
	"flag"
	"github.com/ispeakc0de/load-gen/pkg/log"
	"net/http"
	"os"
	"time"
)

func main() {
	var (
		logFormatter string
		logLevel     string
	)

	flag.StringVar(&logFormatter, "log-formatter", "json", "Log formatter (text|json)")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (trace|debug|info|warn|error|fatal|panic)")
	flag.Parse()

	// Initialize logger
	log.InitLogger(logFormatter, logLevel)

	URL := getenv("URL", "")
	Interval := getenv("INTERVAL", "2s")

	interval, err := time.ParseDuration(Interval)
	if err != nil {
		log.Logger.Fatalf("Failed to parse interval: %v", err)
	}

	log.Logger.Infof("Generating load on URL: %s every %s", URL, Interval)

	httpClient := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for {
		log.Logger.Infof("Sending request to %s", URL)
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			log.Logger.Errorf("Failed to create request: %v", err)
			time.Sleep(interval)
			continue
		}

		_, err = httpClient.Do(req)
		if err != nil {
			log.Logger.Errorf("Failed to send request: %v", err)
			time.Sleep(interval)
			continue
		}
		time.Sleep(interval)
	}
}

func getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
