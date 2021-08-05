package main

import (
	"encoding/json"
    "time"

    "github.com/dhowden/tag"
)

type AudioFile struct {
    filename, artist, albumArtist, album, title, year string
	track, disc										  uint32
    discTracks  									  *map[uint32]uint32
    time        									  time.Duration
    cover       									  *tag.Picture
}

type FormData struct {
    albumDirectory, coverPath   string
    extractCover, separateFiles bool
    outputPath                  string
}

type VideoData struct {
    formData   FormData
    audioFiles []AudioFile
    imageFiles []string
    discTracks map[uint32]uint32
}

type Timestamp struct {
	Artist 		string `json:"artist"`
	AlbumArtist string `json:"albumArtist"`
	Title 		string `json:"title"`
	Time 		string `json:"time"`
}

func (t Timestamp) String() (s string) {
	b, err := json.Marshal(t); if err != nil { panic(err) }
	return string(b)
}
