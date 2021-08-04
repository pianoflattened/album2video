package main

import (
    "encoding/json"
    "errors"
    "os"
    "path"
    "strings"

    "github.com/Akumzy/ipc"
    "github.com/gabriel-vasile/mimetype"
    "github.com/u2takey/ffmpeg-go"
)

func getTags(channel *ipc.IPC, formData FormData) VideoData {
    validatePaths(channel, formData)

    setLabel(channel, "reading " + formData.albumDirectory + "..")

    albumDirectoryFile, err := os.Open(formData.albumDirectory); if err != nil { panic(err) }

    files, err := albumDirectoryFile.Readdirnames(-1); if err != nil { panic(err) }
    
    for _, base := range files {
        file := path.Join(formData.albumDirectory, base)
        mime, err := mimetype.DetectFile(file); if err != nil { panic(err) }

        switch strings.Split(mime.String(), "/")[0] {
        case "audio":
            ffprobeJSON, err := ffmpeg.Probe(file); if err != nil { panic(err) }
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
