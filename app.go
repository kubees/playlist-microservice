package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/go-redis/redis/v9"
	"github.com/julienschmidt/httprouter"
	"github.com/kubees/playlist-microservice/jaeger"
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
var traceProvider = jaeger.NewJaegerTracerProvider()

func main() {

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(traceProvider)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)
	defer Logger.Sync()

	r := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{redisHost + ":" + redisPort},
		DB:       0,
		Password: password,
	})
	rdb = r

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
