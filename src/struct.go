package main

import (
    "time"

    "github.com/dhowden/tag"
)

type AudioFile struct {
    filename    string
    artist      string
    albumArtist string
    album       string
    title       string
    year        string
    track       uint64
    disc        uint64
    discTracks  *map[uint64]uint64
    time        time.Duration
    cover       *tag.Picture
}

type FormData struct {
    albumDirectory string
    coverPath      string
    detectCover    bool
    separateFiles  bool
    outputPath     string
    ffprobePath    string
    ffmpegPath     string
}

type VideoData struct {
    formData   FormData
    audioFiles []AudioFile
    imageFiles []string
    discTracks map[uint64]uint64
}
