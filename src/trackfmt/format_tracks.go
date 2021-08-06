package main

import (
	"fmt"
	"github.com/Akumzy/ipc"
)

func formatTracks(channel *ipc.IPC, fmtString string, tracks string) (o string) {
	/* for _, track := range tracks {
		buf := []rune{}
		o = ""
		for _, r := range fmtString {
			if r == '%' || len(buf) == 1 {
				buf = append(buf, r)
			}
			
			if len(buf) == 2 {
				switch buf[1] {
				case 't':
					o += track.Title
				case 's':
				case 'r':
				case 'a':
				case 'd':
				case 'n':
				case 'w':
				case '%':
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				}
			}
			
			if r == ' ' {
				buf = []rune{}
			}
		}
	} */
	return "lol"
}

func Println(ipc *ipc.IPC, msg interface{}) {
    ipc.Send("log", fmt.Sprintf("%v", msg))
}