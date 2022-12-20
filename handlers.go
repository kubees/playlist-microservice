package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/attribute"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tr := traceProvider.Tracer("playlist-ms-main-component")
	id := uuid.New()
	ip := strings.Split(r.RemoteAddr, ":")[0]
	_, span := tr.Start(context.Background(), "healthz")
	span.SetAttributes(attribute.Key("Protocol").String(r.Proto))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(ip))
	defer span.End()
	Sugar.Infof("client_ip: %v", ip)
	Sugar.Infof("request_id: %v", id.String())
	fmt.Fprintf(w, "ok!")
}

func GetPlaylistsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tr := traceProvider.Tracer("playlist-ms-main-component")
	ctx, span := tr.Start(context.Background(), "GET Playlist HTTP Handler")
	id := uuid.New()
	ip := strings.Split(r.RemoteAddr, ":")[0]
	Sugar.Infof("client_ip: %v", ip)
	Sugar.Infof("request_id: %v", id.String())

	defer span.End()

	span.SetAttributes(attribute.Key("Function").String("GetPlaylistHandler"))
	span.SetAttributes(attribute.Key("Protocol").String(r.Proto))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(strings.Split(r.RemoteAddr, ":")[0]))

	Cors(w)
	playlistsJson := GetPlaylists(ctx, id, ip)

	var playlists []Playlist
	err := json.Unmarshal([]byte(playlistsJson), &playlists)
	if err != nil {
		Sugar.Errorf("Error while unmarshalling JSON file: %v\n", err)
		return
	}
	PlaylistsMetrics(playlists)

	//get videos for each playlist from videos api
	for pi := range playlists {
		vs := GetVideosOfPlaylists(playlists[pi], ctx, id, ip)
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
