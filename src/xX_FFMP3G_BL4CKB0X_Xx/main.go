package main

import (
	"os"

	"github.com/Akumzy/ipc"
)

var ipcIO *ipc.IPC

func main() {
	ipcIO = ipc.New()

	go func() {
		formData := FormData{
			albumDirectory: os.Args[1],
			coverPath:      os.Args[2],
			extractCover:   os.Args[3] == "true",
			separateFiles:  os.Args[4] == "true",
			outputPath:     os.Args[5],
		}
		ffprobePath := os.Args[6]
		ffmpegPath := os.Args[7]

		videoData := getTags(ipcIO, formData, ffprobePath)
		concatWav := makeVideo(ipcIO, videoData, ffmpegPath)
		defer func() {
			err := os.Remove(concatWav)
			if err != nil {
				Println(ipcIO, err)
			}
		}()
	}()

	ipcIO.Start()
}
