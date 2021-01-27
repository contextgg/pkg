package httputils

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
	healthAddr  = ":8082"
)

func Serve(mux http.Handler) {
	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	h := std.Handler("", mdlw, mux)

	// Serve our handler.
	go func() {
		log.Println("Running server")
		if err := http.ListenAndServe(srvAddr, h); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Println("Running metrics")
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	// Serve our health.
	go func() {
		log.Println("Running health check")
		if err := http.ListenAndServe(healthAddr, health); err != nil {
			log.Panicf("error while serving health: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
