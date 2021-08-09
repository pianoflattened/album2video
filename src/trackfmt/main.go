package main

import (
	"fmt"
	"os"
)

func main() {
	fmtStr := os.Args[1]
	timestamps := os.Args[2]

	if len(fmtStr) == 0 {
		fmtStr = "%[ - >a%s - %t"
	}

	formatted := formatTracks(fmtStr, timestamps)
	fmt.Println(formatted)
}

/*'%v[ - >a%vs - %3vt' '[{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Sonata for violin and piano","time":"0:00","disc":1,"track":1,"overallTrack":1},{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Music for violin cello and piano","time":"29:00","disc":1,"track":2,"overallTrack":2},{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Dotted music n°1","time":"1:03:00","disc":1,"track":3,"overallTrack":3}]'*/
