package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetPlaylists() (response string) {
	playlistData, err := rdb.Get(ctx, "playlists").Result()

	if err != nil {
		fmt.Println(err)
		fmt.Println("error occured retrieving playlists from Redis")
		return "[]"
	}

	return playlistData
}

func GetVideosOfPlaylists(playlist Playlist) []Videos {
	vs := []Videos{}
	for vi := range playlist.Videos {

		v := Videos{}
		videoResp, err := http.Get(fmt.Sprintf("http://%v:%v/", videos_api_host, videos_api_port) + playlist.Videos[vi].Id)

		if err != nil {
			fmt.Println(err)
			break
		}

		defer videoResp.Body.Close()
		video, err := ioutil.ReadAll(videoResp.Body)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(video, &v)

		if err != nil {
			panic(err)
		}

		vs = append(vs, v)

	}
	return vs

}