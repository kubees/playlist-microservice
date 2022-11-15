package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var environment = os.Getenv("ENVIRONMENT")
var redisHost = os.Getenv("REDIS_HOST")
var redisPort = os.Getenv("REDIS_PORT")
var videosApiHost = os.Getenv("VIDEOS_API_HOST")
var videosApiPort = os.Getenv("VIDEOS_API_PORT")
var ctx = context.Background()
var rdb *redis.Client

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
		DB:   0,
	})
	rdb = r

	router := httprouter.New()

	router.GET("/healthz", HealthzHandler)

	router.GET("/", GetPlaylistsHandler)

	fmt.Println("Running...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":10010", router))
}
