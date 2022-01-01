package main

import (
	"encoding/json"
	"time"

	"github.com/dhowden/tag"
)

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
	upper
)

type AudioFile struct {
	filename, artist, albumArtist, album, title, year string
	track, disc                                       int
	discTracks                                        *map[int]int
	time                                              time.Duration
	cover                                             *tag.Picture
}

type FormData struct {
	albumDirectory, coverPath   string
	extractCover, separateFiles bool
	outputPath                  string
	verbose, quiet				bool
}

type VideoData struct {
	formData   FormData
	audioFiles []AudioFile
	imageFiles []string
	discTracks map[int]int
}

type Timestamp struct {
	Artist       string `json:"artist"`
	AlbumArtist  string `json:"albumArtist"`
	Title        string `json:"title"`
	Time         string `json:"time"`
	Disc         int    `json:"disc"`
	Track        int    `json:"track"`
	OverallTrack int    `json:"overallTrack"`
}

func (t Timestamp) String() (s string) {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FmtTimestamp struct {
	r, a, s, t string
	d, n, w    int
}

func (t Timestamp) toFmtTimestamp() FmtTimestamp {
	return FmtTimestamp{t.Artist, t.Artist, t.Title, t.Time, t.Disc, t.Track, t.OverallTrack}
}

type renderOptions struct {
	fCase          FieldCase
	padding        int
	ifExistsRight  bool
	ifExists       string
	mode           rune
	dominantArtist string
	multDiscs      bool
}
