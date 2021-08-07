package main

import (
	"encoding/json"
	"strconv"

	"github.com/Akumzy/ipc"
)

/*
[{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Sonata for violin and piano","time":"0:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Music for violin cello and piano","time":"29:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Dotted music n┬░1","time":"1:03:00"}]
*/

func _formatTracks(channel *ipc.IPC, fmtString string, tracks string) string {
	var timestamps []Timestamp
	json.Unmarshal([]byte(tracks), &timestamps)
	o := ""

	dominantArtist, multDiscs := calcDominantArtist(channel, timestamps)
	for i := range timestamps {
		Println(channel, i)
	}

	for _, timestamp := range timestamps {
		Println(channel, timestamp)
		s := fmtState{
			tokenBuf:       []rune{},
			padding:        3,
			fCase:          raw,
			bracketBuf:     []rune{},
			inBrackets:     false,
			bracketEscaped: false,
			renderRight:    false,
			multDiscs:      multDiscs,
		}
		line := ""

		for i, c := range fmtString {
			if s.inBrackets {
				if s.bracketEscaped {
					s.appendB(c)
					s.bracketEscaped = false
					continue
				}

				if c == '\\' {
					s.bracketEscaped = true
					continue
				}

				if (c == '}' && s.renderRight) || (c == '{' && !s.renderRight) {
					s.inBrackets = false
					continue
				}

				if i == len(fmtString)-1 {
					line += string(s.bracketBuf)
					break
				}

				continue
			}

			if s.length() == 1 && c == '%' { // escaped percent char
				line += "%"
				s.clear()
				continue
			}

			if s.length() > 0 {
				Println(channel, string(s.tokenBuf))
				switch c {
				case 't':
					line += render(channel, &s, timestamp.Title, "title")
					s.clear()
				case 's':
					line += render(channel, &s, timestamp.Time, "time")
					s.clear()
				case 'r':
					line += render(channel, &s, timestamp.Artist, "artist")
					s.clear()
				case 'a':
					line += render(channel, &s, discriminate(timestamp.Artist, dominantArtist), "artist")
					s.clear()
				case 'd':
					line += render(channel, &s, strconv.Itoa(timestamp.Disc), "disc")
					s.clear()
				case 'n':
					line += render(channel, &s, strconv.Itoa(timestamp.OverallTrack), "track")
					s.clear()
				case 'w':
					line += render(channel, &s, strconv.Itoa(timestamp.Track), "track")
					s.clear()
				case 'c':
					s.fCase = lower
					s.append(c)
				case 'C':
					s.fCase = title
					s.append(c)
				case '[':
					s.renderRight = true
					if s.length() == 1 {
						s.inBrackets = true
					}
				case '{':
					s.renderRight = false
					if s.length() == 1 {
						s.inBrackets = true
					}
				default:
					if s.length() == 1 {
						s.clear() // consume the next character bc im evil
					}
				}
				continue
			}

			if c == '%' {
				s.append(c)
				continue
			}
			line += string([]rune{c})
			Println(channel, line)
		}
		o += line + "\n"
		Println(channel, o)
	}

	channel.Send("result", o)
	Println(channel, "result: "+o)
	return o
}
