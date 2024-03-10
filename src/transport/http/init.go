package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
	"website-analyzer/src/configuration"
	"website-analyzer/src/controller"
	"website-analyzer/src/metrics"
	"website-analyzer/src/transport/http/engine"
)

const logPrefixHttpInit = "web-analyzer/src/transport/http/engine/init"

var srvHttp http.Server
var srvMetrics http.Server

func Start(analyzerController controller.AnalyzerControllerInterface) {

	srvHttp = http.Server{
		Addr:         ":" + strconv.Itoa(configuration.AppConfig.App.Port),
		Handler:      engine.NewHttpEngine(analyzerController).Engine(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srvMetrics = http.Server{
		Addr:         ":" + strconv.Itoa(configuration.AppConfig.Metrics.Port),
		Handler:      metrics.NewMetricsEngine().Engine(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("%v,%v,%v", "INFO", logPrefixHttpInit, "starting server on port:"+strconv.Itoa(configuration.AppConfig.App.Port))
		err := srvHttp.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("%v,%v,%v,%v", "ERROR", logPrefixHttpInit, "failed to start http server:", err.Error())
		}
	}()

	go func() {
		log.Printf("%v,%v,%v", "INFO", logPrefixHttpInit, "starting metrics server on port:"+strconv.Itoa(configuration.AppConfig.Metrics.Port))
		err := srvMetrics.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("%v,%v,%v,%v", "ERROR", logPrefixHttpInit, "failed to start metrics server:", err.Error())
		}
	}()
}

func Shutdown(ctx context.Context) {
	if err := srvHttp.Shutdown(ctx); err != nil {
		log.Fatal(logPrefixHttpInit, "failed to shutdown http server: ", err)
	}
	if err := srvMetrics.Shutdown(ctx); err != nil {
		log.Fatal(logPrefixHttpInit, "failed to shutdown metrics web server: ", err)
	}
}
