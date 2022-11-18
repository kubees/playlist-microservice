package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	NumberOfPlaylists = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_present_playlists",
		Help: "This is the total number of playlists in the DB.",
	})
	NumberOfVideos = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "number_of_present_videos",
		Help: "This is the total number of videos in the DB.",
	})
	NumberOfVideosPerPlaylist = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "number_of_videos_per_playlist",
		Help: "This is the number of videos per playlist in the DB.",
	}, []string{"playlist"})
)

func RegisterMetrics() {
	prometheus.Register(NumberOfPlaylists)
	prometheus.Register(NumberOfVideos)
	prometheus.Register(NumberOfVideosPerPlaylist)
}

func PlaylistsMetrics(playlists []Playlist) {
	NumberOfPlaylists.Set(float64(len(playlists)))
	for _, playlist := range playlists {
		NumberOfVideosPerPlaylist.WithLabelValues(playlist.Name).Set(float64(len(playlist.Videos)))
		videosMetrics(playlist.Videos)
	}
}

func videosMetrics(videos []Videos) {
	NumberOfVideos.Set(float64(len(videos)))
}
