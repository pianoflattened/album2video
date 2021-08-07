package main

import (
	"math"

	"github.com/Akumzy/ipc"
)

func calcDominantArtist(channel *ipc.IPC, timestamps []Timestamp) (string, bool) {

	var dominantArtist string
	artists := make(map[string]int)
	discs := 0

	for _, track := range timestamps {
		if _, ok := artists[track.Artist]; ok {
			artists[track.Artist] += 1
		} else {
			artists[track.Artist] = 1
		}

		if track.Disc > discs {
			discs = track.Disc
		}
	}

	variousArtists := len(artists) > 1

	highestTrackCount := -1
	if variousArtists {
		for k, v := range artists {
			if v > highestTrackCount {
				highestTrackCount = v
				dominantArtist = k
			}
		}

		if float64(highestTrackCount) < math.Floor((5./6.)*float64(len(timestamps))) {
			dominantArtist = "" // dominantArtist must take up at least 5/6ths of the tracklist to count
		}
	} else {
		dominantArtist = timestamps[0].Artist
	}
	return dominantArtist, discs > 1
}

func discriminate(artist, dominantArtist string) string {
	if dominantArtist != "" {
		if artist != dominantArtist {
			return artist
		}
		return ""
	}
	return artist
}
