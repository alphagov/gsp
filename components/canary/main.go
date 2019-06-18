package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var BuildTimestamp string
var BuildTimeUnixFloat float64
var ChartCommitTimeUnixFloat float64

func init() {
	var err error
	BuildTimeUnixFloat, err = strconv.ParseFloat(BuildTimestamp, 64)

	if err != nil {
		panic(err)
	}

	ChartCommitTimeUnixFloat, err = strconv.ParseFloat(
		os.Getenv("CHART_COMMIT_TIMESTAMP"), 64,
	)

	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	buildTimeMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "canary_build_timestamp",
		},
	)
	buildTimeMetric.Set(BuildTimeUnixFloat)
	prometheus.DefaultRegisterer.MustRegister(buildTimeMetric)

	chartCommitTimeMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "canary_chart_commit_timestamp",
		},
	)
	chartCommitTimeMetric.Set(ChartCommitTimeUnixFloat)
	prometheus.DefaultRegisterer.MustRegister(chartCommitTimeMetric)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
