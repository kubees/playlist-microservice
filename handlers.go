package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "ok!")
}

func GetPlaylistsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	Cors(w)
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

}
