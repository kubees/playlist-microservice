package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func GetPlaylists(ctx context.Context, id uuid.UUID, ip string) (response string) {
	tr := traceProvider.Tracer("playlist-ms-main-component")
	ctx, span := tr.Start(ctx, "GET Playlist From DB")
	defer span.End()
	span.SetAttributes(attribute.Key("Function").String("GetPlaylists"))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(ip))
	playlistData, err := rdb.Get(ctx, "playlists").Result()

	if err == redis.Nil {
		Sugar.Infow("there's no playlists right now!")
		return "[]"
	} else if err != nil {
		span.RecordError(err)
		Sugar.Errorw("Error while trying to retrieve playlists data from redis", err)
		return "[]"
	}
	Sugar.Infow("Returning Playlists Successfully", playlistData)
	return playlistData
}

func GetVideosOfPlaylists(playlist Playlist, ctx context.Context, id uuid.UUID, ip string) []Videos {
	tr := traceProvider.Tracer("playlist-ms-main-component")
	ctx, span := tr.Start(ctx, "Fetch Videos from videos ms")
	defer span.End()
	span.SetAttributes(attribute.Key("Function").String("GetVideosOfPlaylist"))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(ip))
	var vs []Videos
	for vi := range playlist.Videos {

		v := Videos{}

		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%v:%v/", videosApiHost, videosApiPort)+playlist.Videos[vi].Id, nil)
		if err != nil {
			Sugar.Errorw("Error while creating a new request with ctx")
			span.RecordError(err)
		}

		videoResp, err := http.DefaultClient.Do(req)

		if err != nil {
			span.RecordError(err)
			Sugar.Errorw("Error while trying to fetch videos from videos microservice", err)
			break
		}

		defer videoResp.Body.Close()
		video, err := ioutil.ReadAll(videoResp.Body)

		if err != nil {
			span.RecordError(err)
			Sugar.Errorw("Error while trying to access video object", err)
			panic(err)
		}

		err = json.Unmarshal(video, &v)

		if err != nil {
			span.RecordError(err)
			Sugar.Errorw("Error while trying to unmarshal video object", err)
			panic(err)
		}

		vs = append(vs, v)

	}
	return vs

}
