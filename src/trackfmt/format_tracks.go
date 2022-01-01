package main

import (
	"crypto/md5"
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kelindar/binary" // sorry
)

var err error // line 58 :(

/*\%(([\^\-v])?(\d+)|(\d+)([\^\-v])?|([\^\-v]))?((\[)(([^\>]|\\.)+)?\}|(\<)(([^\]]|\\.)+)?\])?([tsradnw\%])
   \% - initial percent
   (([\^\-v])?(\d+)|(\d+)([\^\-v])?|([\^\-v]))? - case/padding
   (
	(\[)(([^\>]|\\.)+)?\> - ifexists right
	|
	(\<)(([^\]]|\\.)+)?\] - ifexists left
   )?
   ([tsradnw\%]) - mode
*/

func formatTracks(fmtString string, timestamps []Timestamp) string {
	o := ""

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
		o += line + "\n"
	}
	return o
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
