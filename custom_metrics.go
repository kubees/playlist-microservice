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
)

func PlaylistsMetrics(playlists []Playlist) {
	NumberOfPlaylists.Set(float64(len(playlists)))
}
