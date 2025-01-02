package main

import (
	"net/http"
	"time"

	"3x-ui-monitoring/auth"
	"3x-ui-monitoring/config"
	"3x-ui-monitoring/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.LoadConfig()
	m := metrics.NewMetrics(cfg)
	go m.StartPolling(5 * time.Second)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8080", auth.BasicAuth(mux, cfg))
}
