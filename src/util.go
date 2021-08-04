package main

import "github.com/Akumzy/ipc"

func setLabel(ipc *ipc.IPC, msg string) {
    ipc.Send("progress-label", msg)
}

func Println(ipc *ipc.IPC, msg string) {
    ipc.Send("log", msg)
}
