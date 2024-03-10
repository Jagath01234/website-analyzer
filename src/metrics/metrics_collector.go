package metrics

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestSummary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "http_request_count",
		Help: "Total number of HTTP requests.",
	},
	[]string{"method", "endpoint", "status"},
)
