package main

type Playlist struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Videos []Videos `json:"videos"`
}

type Videos struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageurl"`
	Url         string `json:"url"`
}
