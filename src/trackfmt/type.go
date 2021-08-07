package main

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
)

type JSONTimestamp struct {
	Artist       string `json:"artist"`
	AlbumArtist  string `json:"albumArtist"`
	Title        string `json:"title"`
	Time         string `json:"time"`
	Disc         int    `json:"disc"`
	Track        int    `json:"track"`
	OverallTrack int    `json:"overallTrack"`
}

func (j JSONTimestamp) toTimestamp() Timestamp {
	return Timestamp{j.Artist, j.Artist, j.Title, j.Time, j.Disc, j.Track, j.OverallTrack}
}

type Timestamp struct {
	r, a, s, t string
	d, n, w    int
}
