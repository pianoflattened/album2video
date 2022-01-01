package main

import (
	"crypto/md5"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kelindar/binary" // sorry
)

var err error // line 67 :(

/* \%(([\^\-v])?(\d+)|(\d+)([\^\-v])?|([\^\-v]))?((\[)(([^\>]|\\.)+)?\}|(\<)(([^\]]|\\.)+)?\])?([tsradnw\%])
 * \% - initial percent
 * (([\^\-v])?(\d+)|(\d+)([\^\-v])?|([\^\-v]))? - case/padding
 * (
 * (\[)(([^\>]|\\.)+)?\> - ifexists right
 * |
 * (\<)(([^\]]|\\.)+)?\] - ifexists left
 * )?
 * ([tsradnw\%]) - mode
*/

func formatTimestamps(fmtString string, timestamps []Timestamp) (formatted string) {
	if len(fmtString) == 0 {
		fmtString = "%[ - >a%s - %t"
	}
	
	formatted = ""

	dominantArtist, multDiscs := calcDominantArtist(timestamps)

	re := regexp.MustCompile(`\%(([\^\-v])?(\d+)|(\d+)([\^\-v])?|([\^\-v]))?((\[)(([^>]|\\.)+)?\>|(<)(([^\]]|\\.)+)?\])?([tsradnw\%])`)
	matches := removeDuplicates(re.FindAllStringSubmatch(fmtString, -1))

	for _, jtimestamp := range timestamps {
		timestamp := jtimestamp.toFmtTimestamp()
		line := fmtString
		r := renderOptions{raw, -1, false, "", '_', dominantArtist, multDiscs}
		for _, match := range matches {
			field := ""
			if match[14] == "%" {
				continue
			}
			r.mode, _ = utf8.DecodeRuneInString(match[14])

			switch firstNonZero(match[2], match[5], match[6]).(string) {
			case "v":
				r.fCase = lower
			case "-":
				r.fCase = title
			case "^":
				r.fCase = upper
			}

			r.padding, err = strconv.Atoi(firstNonZero(match[3], match[4]).(string))
			if err != nil {
				r.padding = -1
			}
			r.ifExistsRight = firstNonZero(match[8], match[11]).(string) == "["
			r.ifExists = firstNonZero(match[9], match[12]).(string)

			field = r.render(timestamp)
			line = strings.ReplaceAll(line, match[0], field)
		}
		// i think this is the more robust way of doing it? skip over escaped % matches and then go back
		// and fix them all?? may have to correct this
		line = strings.ReplaceAll(line, "%%", "%")
		formatted += line + "\n"
	}
	return
}

func firstNonZero(n interface{}, m ...interface{}) interface{} {
	v := reflect.ValueOf(n)
	if len(m) <= 1 {
		if v.IsZero() {
			return m[0]
		}
		return n
	}
	if v.IsZero() {
		return firstNonZero(m[0], m[1:]...)
	}
	return ""
}

func removeDuplicates(a [][]string) (b [][]string) {
	keys := make(map[string]bool)
	b = [][]string{}
	for _, e := range a {
		g, err := binary.Marshal(e)
		if err != nil {
			panic(err)
		}
		sum := md5.Sum(g) // what the fuck was i on why am i using md5 here. fucking horrifying
		if _, ok := keys[string(sum[:])]; !ok {
			keys[string(sum[:])] = true
			c := make([]string, len(e))
			copy(c, e)
			b = append(b, c)
		}
	}
	return b
}

func calcDominantArtist(timestamps []Timestamp) (string, bool) {
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

func (r renderOptions) render(t FmtTimestamp) string {
	s := r.valFromMode(t)

	if strings.IndexRune("tdnw", r.mode) > -1 { // padding
		if r.mode == 't' {
			s = fmt.Sprintf("%0*s", int(math.Max(float64(r.padding), 3)), strings.ReplaceAll(s, ":", ""))
			if len(s) > 4 {
				s = s[0:len(s)-4] + ":" + s[len(s)-4:]
			}
			s = s[0:len(s)-2] + ":" + s[len(s)-2:]
		} else {
			s = fmt.Sprintf("%0*s", r.padding, s)
		}
	}

	if r.mode == 'd' && !r.multDiscs {
		s = ""
	}

	if r.mode == 'a' {
		s = discriminate(s, r.dominantArtist)
	}

	if s != "" {
		if r.ifExistsRight {
			s = s + r.ifExists
		} else {
			s = r.ifExists + s
		}
	}

	if r.fCase == lower { // case
		s = strings.ToLower(s)
	} else if r.fCase == title {
		s = strings.Title(s) // SUCKS
	}

	return s
}

// time 2 whip out th reflect library for opaque bullshit >:)
func (r renderOptions) valFromMode(t FmtTimestamp) string {
	switch f := reflect.ValueOf(t).FieldByName(string([]rune{r.mode})); f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(f.Int(), 10)
	case reflect.String:
		return f.String()
	}
	return ""
}
