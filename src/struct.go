package main

import (
    "time"

    "github.com/dhowden/tag"
)

type AudioFile struct {
    filename, artist, albumArtist, album, title, year string
	track, disc										  uint64
    discTracks  									  *map[uint64]uint64
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
    discTracks map[uint64]uint64
}

type Timestamp struct {
	artist, title, time string
}
