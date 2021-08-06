package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	
	"github.com/Akumzy/ipc"
)

type Timestamp struct {
	Artist 		 string `json:"artist"`
	AlbumArtist  string `json:"albumArtist"`
	Title 		 string `json:"title"`
	Time 		 string `json:"time"`
	Disc		 int    `json:"disc"`
	Track		 int    `json:"track"`
	OverallTrack int	`json:"overallTrack"`
}

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
)

/*
[{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Sonata for violin and piano","time":"0:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Music for violin cello and piano","time":"29:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Dotted music n┬░1","time":"1:03:00"}]
*/

func formatTracks(channel *ipc.IPC, fmtString string, tracks string) (o string) {
	var timestamps []Timestamp
	json.Unmarshal([]byte(tracks), &timestamps)
	
	var dominantArtist string
	artists 	 	   := make(map[string]int)
	albumArtists 	   := make(map[string]int)
	albumArtistUnknown := true
	
	for _, track := range timestamps {
		if _, ok := artists[track.Artist]; ok {
			artists[track.Artist] += 1
		} else {
			artists[track.Artist] = 1
		}
		
		if _, ok := albumArtists[track.AlbumArtist]; ok {
			albumArtists[track.AlbumArtist] += 1
		} else {
			albumArtists[track.AlbumArtist] = 1
		}
		
		if track.AlbumArtist != "[unknown artist]" { albumArtistUnknown = false }
	}
	
	variousArtists 		:= len(artists) > 1
	variousAlbumArtists := len(albumArtists) > 1
	
	highestTrackCount := -1
	if variousArtists {
		for k, v := range artists {
			if v > highestTrackCount {
				highestTrackCount = v
				dominantArtist = k
			}
		}
		
		if highestTrackCount < math.Floor((5./6.)*len(timestamps)) {
			dominantArtist = "" // dominantArtist must take up at least 5/6ths of the tracklist to count
		}
	} else {
		dominantArtist = timestamps[0].Artist
	}
	
	if dominantArtist != "[unknown artist]" { albumArtistUnknown = false }

	var readingArgs bool
	var fieldCase 	FieldCase
	var padding 	int
	var c 			rune
	var field 		string
	
	for _, track := range timestamps {
		buf := []rune{}
		o = ""
		numberStart = -1
		fieldCase = raw
		padding = 3
		for i, r := range fmtString {
			if len(buf) >= 1 {
				ind := Index(r, []rune("0123456789tsradnw%")) // add braces later
				if 0 <= ind <= 9 {
					if numberStart == -1 {
						numberStart = i
					}
					buf = append(buf, r)
					continue
				} else if numberStart > -1 {
					padding, _ = strconv.Atoi(string(buf[numberStart:]))
				}
				
				if 10 <= ind <= 17 {
					switch r {
					case 't':
						field = track.Title
					case 's':
						field = track.Time
					case 'r':
						field = track.Artist
					case 'a':
						field = track.Artist // implement discriminate artist later
					case 'd':
						field = track.Disc
					case 'n':
						field = track.OverallTrack
					case 'w':
						field = track.Track
					case '%':
						field = "\n" // wacky stuff
					}
					o += render(field, fieldCase, padding)
					continue
				}
				
				if r == 'c' {
					fieldCase = lower
					continue
				} else if r == 'C' {
					fieldCase = title
					continue
				}
			}
			
			if len(buf) == 0 && r == '%' {
				buf = append(buf, r)
			}

			if len(buf) >= 2 {
				
				switch buf[len(buf)-1] {
				case 't':
					
				case 's':
				case 'r':
				case 'a':
				case 'd':
				case 'n':
				case 'w':
				case '%':
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					readingNumber = true
				case 'c':
					charCase = "lower"
				case 'C':
					charCase = "title"
				case 'r':
					charCase = "raw"
				}
				buf = append(buf, r)
			}
			
			if r == ' ' {
				buf = []rune{}
			}
		}
	}
	return "lol"
}

func render(field string, fieldCase FieldCase, padding int) string {
	if padding > 3 {
		// lol not doing this today
	}

	switch FieldCase {
	case raw:
		return field
	case lower:
		return strings.ToLower(field)
	case title:
		return field // also not doing this >:)
	}
}

func Println(ipc *ipc.IPC, msg interface{}) {
    ipc.Send("log", fmt.Sprintf("%v", msg))
}

func Index(r rune, buf []rune) int {
	for i, e := range buf {
		if e == r {
			return i
		}
	}
	return -1
}