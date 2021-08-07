package main

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
)

type Timestamp struct {
	Artist       string `json:"artist"`
	AlbumArtist  string `json:"albumArtist"`
	Title        string `json:"title"`
	Time         string `json:"time"`
	Disc         int    `json:"disc"`
	Track        int    `json:"track"`
	OverallTrack int    `json:"overallTrack"`
}
