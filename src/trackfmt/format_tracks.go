package main

import (
	"encoding/json"

	"github.com/Akumzy/ipc"
)

/*
[{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Sonata for violin and piano","time":"0:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Music for violin cello and piano","time":"29:00"},
{"artist":"Taku Sugimoto","albumArtist":"[unknown artist]","title":"Dotted music n┬░1","time":"1:03:00"}]
*/

func formatTracks(channel *ipc.IPC, fmtString string, tracks string) (o string) {
	var timestamps []Timestamp
	json.Unmarshal([]byte(tracks), &timestamps)

	dominantArtist, multDiscs := calcDominantArtist(channel, timestamps)

	for _, timestamp := range timestamps {
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
				switch c {
				case 't':
					line += s.render(timestamp.Title) // render auto-clears
				case 's':
					line += s.render(timestamp.Time, true)
				case 'r':
					line += s.render(timestamp.Artist)
				case 'a':
					line += s.render(discriminate(timestamp.Artist, dominantArtist))
				case 'd':
					line += s.render(timestamp.Disc, false)
				case 'n':
					line += s.render(timestamp.OverallTrack)
				case 'w':
					line += s.render(timestamp.Track)
				case '[':
					s.renderRight = true
					if s.length() == 1 {
						s.inBrackets = true
						continue
					}
				case '{':
					s.renderRight = false
					if s.length() == 1 {
						s.inBrackets = true
						continue
					}
				default:
					if s.length() == 1 {
						s.clear() // consume the next character bc im evil
						continue
					}
				}
			}

			if c == '%' {
				s.append(c)
				continue
			}

			line += string([]rune{c})
		}
		o += line + "\n"
	}
	return
}
