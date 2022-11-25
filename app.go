package main

import (
	"go.uber.org/zap"
	"net/http"
	"os"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
)

const metricsAddr = ":8000"

var environment = os.Getenv("ENVIRONMENT")
var redisHost = os.Getenv("REDIS_HOST")
var redisPort = os.Getenv("REDIS_PORT")
var videosApiHost = os.Getenv("VIDEOS_API_HOST")
var videosApiPort = os.Getenv("VIDEOS_API_PORT")
var password = os.Getenv("PASSWORD")
var rdb redis.UniversalClient
var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()

func main() {
	defer Logger.Sync()

	r := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{redisHost + ":" + redisPort},
		DB:       0,
		Password: password,
	})
	rdb = r
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(r); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(r); err != nil {
		panic(err)
	}
	RegisterMetrics()

	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	router := httprouter.New()
	router.GET("/healthz", httproutermiddleware.Handler("/healthz", HealthzHandler, mdlw))
	router.GET("/", httproutermiddleware.Handler("/", GetPlaylistsHandler, mdlw))
	Sugar.Infof("Running...")
	// Serve our metrics.
	go func() {
		Sugar.Infof("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			Sugar.Panicf("error while serving metrics: %s", err)
		}
	}()

	log.Fatal(http.ListenAndServe(":10010", router))

}
