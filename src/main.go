package main

import (
    "os"

    "github.com/Akumzy/ipc"
)

var ipcIO *ipc.IPC

type Who struct {
    Name string `json:"name,omitempty"`
}

// i wrote out like three whole long ass comments about how much i hated this one github user for
// making what i think was a mimetype detection library laughably unusable so i guess to counteract
// some of that negative energy shout out to akumzy i love you so much dude you made a module for
// exactly what i wanted to accomplish (ipc between a golang binary and node js) and the examples
// were clear as hell and made me laugh you did such an incredible job and youve made my life so
// much easier if yr reading this then thank u so much :DDDDD
func main() {
    ipcIO = ipc.New()

    go func() {
        formData := FormData{
            albumDirectory: os.Args[1],
            coverPath: os.Args[2],
            detectCover: os.Args[3] == "true",
            separateFiles: os.Args[4] == "true",
            outputPath: os.Args[5],
            ffprobePath: os.Args[6],
            ffmpegPath: os.Args[7],
        }
        ffprobePath := os.Args[6]
        //ffmpegPath := os.Args[7]

        videoData := getTags(ipcIO, formData, ffprobePath)
        makeVideo(ipcIO, videoData)
    }()

    ipcIO.Start()
}
