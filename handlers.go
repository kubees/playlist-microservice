package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/attribute"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tr := traceProvider.Tracer("playlist-ms-main-component")

	_, span := tr.Start(context.Background(), "healthz")
	defer span.End()
	fmt.Fprintf(w, "ok!")
}

func GetPlaylistsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tr := traceProvider.Tracer("playlist-ms-main-component")
	ctx, span := tr.Start(context.Background(), "GET Playlist")
	defer span.End()
	span.SetAttributes(attribute.Key("Function").String("GetPlaylistHandler"))

	Cors(w)
	playlistsJson := GetPlaylists(ctx)

	var playlists []Playlist
	err := json.Unmarshal([]byte(playlistsJson), &playlists)
	if err != nil {
		Sugar.Errorf("Error while unmarshalling JSON file: %v\n", err)
		return
	}
	PlaylistsMetrics(playlists)

	//get videos for each playlist from videos api
	for pi := range playlists {
		vs := GetVideosOfPlaylists(playlists[pi], ctx)
		playlists[pi].Videos = vs
	}

	playlistsBytes, err := json.Marshal(playlists)
	if err != nil {
		Sugar.Errorf("Error while marshalling JSON file: %v\n", err)
		return
	}

	reader := bytes.NewReader(playlistsBytes)
	if b, err := ioutil.ReadAll(reader); err == nil {
		fmt.Fprintf(w, "%s", string(b))
	} else {
		Sugar.Errorf("Error while reading data: %v\n", err)
	}

}
