package main

import (
    //"io/ioutil"
    //"os"
    //"path"
    "time"

	"github.com/Akumzy/ipc"
    ffmpeg "github.com/modfy/fluent-ffmpeg"
)

func makeVideo(channel *ipc.IPC, videoData VideoData, ffmpegPath string, ffprobePath string) {
	timestamps := []Timestamp{}
	for f := range videoData.audioFiles {
		timestamps = append(Timestamp{
			artist: f.artist,
			title: f.title,
			time: f.time,
		})
	}

    ffmpeg.SetFfProbePath(ffprobePath)
   	stage1 := ffmpeg.NewCommand(ffmpegPath)
}
