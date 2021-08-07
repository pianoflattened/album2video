package main

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/Akumzy/ipc"
)

/*
%t title
%s timestamp
%r artist (indiscriminate)
%a artist (discriminate)
%d disc
%n track number (overall)
%w track number (within disc)
%% percent
*/

func formatTracks(channel *ipc.IPC, fmtString string, tracks string) string {
	var timestamps []Timestamp
	json.Unmarshal([]byte(tracks), &timestamps)
	o := ""

	dominantArtist, _ := calcDominantArtist(channel, timestamps)
	for i := range timestamps {
		Println(channel, i)
	}

	for _, timestamp := range timestamps {
		Println(channel, timestamp)
		line := fmtString
		line = replaceVal(line, timestamp.Title, "t")
		line = replaceVal(line, timestamp.Time, "s")
		line = replaceVal(line, timestamp.Artist, "r")
		line = replaceVal(line, discriminate(timestamp.Artist, dominantArtist), "a")
		line = replaceVal(line, strconv.Itoa(timestamp.Disc), "d")
		line = replaceVal(line, strconv.Itoa(timestamp.OverallTrack), "n")
		line = replaceVal(line, strconv.Itoa(timestamp.Track), "w")
		strings.ReplaceAll(string(line), "%%", "%")
		o += string(line) + "\n"
	}
	Println(channel, o)
	return o
}

func replaceVal(line, val, letter string) string {
	re := regexp.MustCompile(`(^|[^%])(\%` + letter + `)`)
	return string(re.ReplaceAll([]byte(line), []byte("${1}"+val)))
}
