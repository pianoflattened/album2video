package main

import (
    "encoding/json"
    "errors"
    "os"
    "path"
    "regexp"    
    "strings"

    "github.com/Akumzy/ipc"
    "github.com/gabriel-vasile/mimetype"
    "github.com/tidwall/gjson"
    "github.com/u2takey/ffmpeg-go"
)

func getTags(channel *ipc.IPC, formData FormData) VideoData {
    validatePaths(channel, formData)

    setLabel(channel, "reading " + formData.albumDirectory + "..")

    albumDirectoryFile, err := os.Open(formData.albumDirectory); if err != nil { panic(err) }

    files, err := albumDirectoryFile.Readdirnames(-1); if err != nil { panic(err) }
    
    discTracks := make(map[int]int)
    audioFiles := []AudioFile
    imageFiles := []string

    // wrote this myself :D i'll probably have to change it sometime
    trackRe := regexp.MustCompile(`^([0-9]+|[A-Za-z]|[0-9]+[A-Za-z]|)(-| - |_| |)([0-9]+|[A-Za-z])(?=. | |_)`)

    for _, base := range files {
        file := path.Join(formData.albumDirectory, base)
        mime, err := mimetype.DetectFile(file); if err != nil { panic(err) }

        switch strings.Split(mime.String(), "/")[0] {
        case "audio": // https://github.com/u2takey/ffmpeg-go#show-ffmpeg-progress
            var disc, track uint64
            var artist, albumArtist, title string
            ffprobeJSON, err := ffmpeg.Probe(file); if err != nil { panic(err) }
            
            metadata := gjson.Get(ffprobeJSON, "format").Raw
            missing := checkMissingMetadata(metadata)
            for _, e := range missing {
                if e == "filename" {
                    filename = file
                }
            }

            if gjson.get(metadata, "tags.disc").Exists() {
                disc = parseTrack(metadata, "tags.disc")
            } else { disc = 1 }
            if gjson.get(metadata, "tags.track").Exists() {
                track = parseTrack(metadata, "tags.track")
            } else {
                if !trackRe.MatchString(file) {
                    panic(errors.New("please make sure your filenames start with a track number if they are not tagged properly (which would be preferrable). for exact specifications as to what does and does not get detected as a track number see https://github.com/sunglasseds/album2video"))
                }
            }
        }
    }

    return VideoData{}
}

func validatePaths(channel *ipc.IPC, formData FormData) {
    setLabel(channel, "validating album path..")
    stats, err := os.Stat(formData.albumDirectory); if err != nil { panic(err) }

    if stats.Mode().IsRegular() {
        formData.albumDirectory = path.Dir(formData.albumDirectory)
    } else if stats.IsDir() {
        // pass
    } else {
        panic(errors.New("album directory is not a file or directory"))
    }

    if !formData.detectCover {
        setLabel(channel, "validating cover path..")
        stats, err = os.Stat(formData.coverPath); if err != nil { panic(err) }

        if stats.Mode().IsRegular() {
            formData.albumDirectory = path.Dir(formData.albumDirectory)
        } else {
            panic(errors.New(formData.coverPath + " is not a file"))
        }
    }

    setLabel(channel, "validating output path..")
    stats, err = os.Stat(formData.outputPath); if err != nil { panic(err) }

    if formData.separateFiles && stats.Mode().IsRegular() {
        formData.outputPath = path.Dir(formData.outputPath)
    } else if !formData.separateFiles && stats.IsDir() {
        formData.outputPath = path.Join(formData.outputPath, "out.mp4")
    } else if formData.separateFiles == stats.IsDir() {
        // pass
    } else {
        panic(errors.New(formData.outputPath + " is not a file or directory"))
    }
}

func checkMissingMetadata(metadata string) (missing []string) {
    fields := []string{"tags.artist", "tags.albumArtist", "tags.title"}
    for _, field := range fields {
        if !gjson.Get(metadata, field).Exists() {
            missing = append(missing, field)
        }
    }
    return
}

func parseMetadata(metadata string) (audioFile AudioFile) {
    
}

func parseTrack(m string, id string) (t uint) {
    t, err = strconv.ParseUint(strings.Split(gjson.get(m, id).String(), "/")[0])
    if err != nil { panic(err) }; return
}
