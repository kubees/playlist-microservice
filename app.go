package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var environment = os.Getenv("ENVIRONMENT")
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var videos_api_host = os.Getenv("VIDEOS_API_HOST")
var videos_api_port = os.Getenv("VIDEOS_API_PORT")
var ctx = context.Background()
var rdb *redis.Client

func main() {

	router := httprouter.New()

	router.GET("/healthz", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, "ok")
	})

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cors(w)
		playlistsJson := GetPlaylists()

		playlists := []Playlist{}
		err := json.Unmarshal([]byte(playlistsJson), &playlists)
		if err != nil {
			panic(err)
		}

		//get videos for each playlist from videos api
		for pi := range playlists {
			vs := GetVideosOfPlaylists(playlists[pi])
			playlists[pi].Videos = vs
		}

		playlistsBytes, err := json.Marshal(playlists)
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(playlistsBytes)
		if b, err := ioutil.ReadAll(reader); err == nil {
			fmt.Fprintf(w, "%s", string(b))
		}

	})

	r := redis.NewClient(&redis.Options{
		Addr: redis_host + ":" + redis_port,
		DB:   0,
	})
	rdb = r

	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":10010", router))
}
