package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Akumzy/ipc"
)

func setLabel(ipc *ipc.IPC, msg string) {
	ipc.Send("progress-label", msg)
}

func setProgress(ipc *ipc.IPC, progress float32) {
	ipc.Send("set-progress", progress)
}

func makeDeterminate(ipc *ipc.IPC) {
	ipc.Send("make-determinate", "GO!!!!!!!!")
}

func sendTimestamps(ipc *ipc.IPC, timestamps []Timestamp) {
	f, err := os.Create("C:\\Users\\user\\Documents\\timestamps.txt")
	if err != nil {
		panic(err)
	}
	n, err := json.Marshal(timestamps)
	if err != nil {
		panic(err)
	}
	f.Write(n)

	ipc.Send("timestamps", timestamps)
}

func Println(ipc *ipc.IPC, msg interface{}) {
	ipc.Send("log", fmt.Sprintf("%v", msg))
}

func setComplete(ipc *ipc.IPC) {
	ipc.Send("set-complete", "AAAAAAAAAA")
}
