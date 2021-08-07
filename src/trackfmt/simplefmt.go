package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"

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
	var timestamps []JSONTimestamp
	json.Unmarshal([]byte(tracks), &timestamps)
	o := ""

	dominantArtist, _ := calcDominantArtist(channel, timestamps)
	for i := range timestamps {
		Println(channel, i)
	}

	/*
	   \% - initial percent
	   (([cC])?(\d+)|(\d+)([cC])?|([cC]))? - case/padding
	   (
	   	(\[)(([^}]|\\.)+)?\} - ifexists right
	   	|
	   	(\{)(([^}]|\\.)+)?\] - ifexists left
	   )?
	   ([tsradnw\%]) - character

	*/

	// wacky stuff
	re := regexp.MustCompile(`\%(([cC])?(\d+)|(\d+)([cC])?|([cC]))?((\[)(([^}]|\\.)+)?\}|(\{)(([^}]|\\.)+)?\])?([tsradnw\%])`)
	matches := re.FindAllStringSubmatch(fmtString, -1)

	for _, jtimestamp := range timestamps {
		timestamp := jtimestamp.toTimestamp()
		line := ""
		r := renderOptions{raw, -1, false, "", '_'}
		for _, match := range matches {
			if match[14] == "%" {
				line += match[14]
				continue
			}
			r.mode, _ = utf8.DecodeRuneInString(match[14])

			switch firstNonZero(match[2], match[5], match[6]).(string) {
			case "c":
				r.fCase = lower
			case "C":
				r.fCase = title
			}

			r.padding, _ = firstNonZero(match[3], match[4]).(int)
			r.ifExistsRight = firstNonZero(match[8], match[11]).(string) == "["
			r.ifExists = firstNonZero(match[9], match[12]).(string)

			line += r.render(timestamp)
		}
		Println(channel, timestamp)
		o += line + "\n"
	}
	Println(channel, o)
	return o
}

func firstNonZero(n interface{}, m ...interface{}) interface{} {
	v := reflect.ValueOf(n)
	if len(m) <= 1 {
		if v.IsZero() {
			return fmt.Sprintf("%v", m)
		}
		return fmt.Sprintf("%v", n)
	}
	if v.IsZero() {
		return firstNonZero(m[0], m[1:])
	}
	return ""
}
