package main

import (
	//"encoding/json"
	"fmt"
	"os"
	//"strings"

	"github.com/Akumzy/ipc"
)

var ipcIO *ipc.IPC

// ipc is so ridiculously easy with this :DDD
func main() {
	ipcIO = ipc.New()

	go func() {
		fmtStr := os.Args[1]
		timestamps := os.Args[2]
		Println(ipcIO, fmtStr+" | "+timestamps)

		resp := formatTracks(ipcIO, fmtStr, timestamps)

		ipcIO.Send("parsed-titles", resp)
	}()

	ipcIO.Start()
}

func Println(ipc *ipc.IPC, msg interface{}) {
	ipc.Send("log", fmt.Sprintf("%v", msg))
}
