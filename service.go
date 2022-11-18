package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"net/http"
)

func GetPlaylists(ctx context.Context) (response string) {
	playlistData, err := rdb.Get(ctx, "playlists").Result()

	if err == redis.Nil {
		fmt.Println("there's no playlists right now!")
		return "[]"
	} else if err != nil {
		fmt.Println(err)
		return "[]"
	}
	return playlistData
}

func GetVideosOfPlaylists(playlist Playlist) []Videos {
	var vs []Videos
	for vi := range playlist.Videos {

		v := Videos{}
		videoResp, err := http.Get(fmt.Sprintf("http://%v:%v/", videosApiHost, videosApiPort) + playlist.Videos[vi].Id)

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
