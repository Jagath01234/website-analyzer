package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsEngine struct {
}

func NewMetricsEngine() *MetricsEngine {
	return &MetricsEngine{}
}

func (e MetricsEngine) Engine() *gin.Engine {
	engine := gin.Default()

	prometheus.MustRegister(HttpRequestSummary)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return engine
}
