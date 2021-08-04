package main

import (
    "strconv"
    "strings"

    "github.com/Akumzy/ipc"
    "github.com/tidwall/gjson"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func setLabel(ipc *ipc.IPC, msg string) {
    ipc.Send("progress-label", msg)
}

func Println(ipc *ipc.IPC, msg string) {
    ipc.Send("log", msg)
}
