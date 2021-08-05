package main

import (
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
    albumDirectory, coverPath           string
    detectCover, separateFiles          bool
    outputPath, ffprobePath, ffmpegPath string
}

type VideoData struct {
    formData   FormData
    audioFiles []AudioFile
    imageFiles []string
    discTracks map[uint32]uint32
}

type Timestamp struct {
	artist, title, time string
}
