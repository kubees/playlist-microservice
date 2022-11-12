package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var environment = os.Getenv("ENVIRONMENT")
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var videos_api_host = os.Getenv("VIDEOS_API_HOST")
var videos_api_port = os.Getenv("VIDEOS_API_PORT")
var ctx = context.Background()
var rdb *redis.Client

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: redis_host + ":" + redis_port,
		DB:   0,
	})
	rdb = r

	router := httprouter.New()

	router.GET("/healthz", HealthzHandler)

	router.GET("/", GetPlaylistsHandler)


	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":10010", router))
}
