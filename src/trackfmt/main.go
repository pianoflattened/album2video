package main

import (
	//"encoding/json"
	"fmt"
	"os"
	//"strings"
	//"github.com/Akumzy/ipc"
)

//var ipcIO *ipc.IPC

// ipc is so ridiculously easy with this :DDD
func main() {
	//ipcIO = ipc.New()

	//finished := make(chan bool)
	//go func() {
	fmtStr := os.Args[1]
	timestamps := os.Args[2]

	formatted := formatTracks( /*ipcIO,*/ fmtStr, timestamps)
	//ipcIO.Send("result", formatted)
	fmt.Println(formatted)
	//finished <- true
	//}()

	//ipcIO.Start()
	//<-finished
	//os.Exit(0)
}

//func Println(ipc *ipc.IPC, msg interface{}) {
//	ipc.Send("log", fmt.Sprintf("%v SYS:%#v", msg, msg))
//}

/*'%c[ - }a%cs - %3ct' '[{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Sonata for violin and piano","time":"0:00","disc":1,"track":1,"overallTrack":1},{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Music for violin cello and piano","time":"29:00","disc":1,"track":2,"overallTrack":2},{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Dotted music n°1","time":"1:03:00","disc":1,"track":3,"overallTrack":3}]'*/
