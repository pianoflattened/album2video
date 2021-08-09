package main

import (
	"fmt"

	"github.com/Akumzy/ipc"
)

func setLabel(ipc *ipc.IPC, msg string) {
	ipc.Send("progress-label", msg)
}

func setProgress(ipc *ipc.IPC, progress float64) {
	ipc.Send("set-progress", progress)
}

func makeDeterminate(ipc *ipc.IPC) {
	ipc.Send("make-determinate", "GO!!!!!!!!")
}

func sendTimestamps(ipc *ipc.IPC, timestamps []Timestamp) {
	ipc.Send("timestamps", timestamps)
}

func Println(ipc *ipc.IPC, msg interface{}) {
	ipc.Send("log", fmt.Sprintf("%v", msg))
}

func setComplete(ipc *ipc.IPC) {
	ipc.Send("set-complete", "AAAAAAAAAA")
}
