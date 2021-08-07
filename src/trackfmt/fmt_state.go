package main

import (
	"math"
	"strings"

	"github.com/Akumzy/ipc"
)

type fmtState struct {
	tokenBuf       []rune
	padding        int
	fCase          FieldCase
	bracketBuf     []rune
	renderRight    bool
	bracketEscaped bool
	multDiscs      bool
	inBrackets     bool
}

func (s *fmtState) append(r rune) {
	(*s).tokenBuf = append((*s).tokenBuf, r)
}

func (s fmtState) length() int {
	return len(s.tokenBuf)
}

func (s *fmtState) appendB(r rune) {
	(*s).bracketBuf = append((*s).bracketBuf, r)
}

func (s *fmtState) clear() {
	(*s).tokenBuf = []rune{}
	(*s).padding = -1
	(*s).fCase = raw
	(*s).inBrackets = false
	(*s).bracketEscaped = false
	(*s).bracketBuf = []rune{}
}

func render(channel *ipc.IPC, s *fmtState, n, field string) string {
	Println(channel, []string{n, field})
	switch field {
	case "artist":
		fallthrough
	case "title":
		if (*s).fCase == lower {
			n = strings.ToLower(field)
		} else if (*s).fCase == title {
			n = TitleCase(n, false)
		}
	case "time":
		n = pad(strings.Join(strings.Split(n, ":"), ""), int(math.Max(float64((*s).padding), 3)))
		if len(n) > 4 {
			n = n[0:len(n)-4] + ":" + n[len(n)-4:]
		}
		n = n[0:len(n)-2] + ":" + n[len(n)-2:]
	case "disc":
		fallthrough
	case "track":
		p := (*s).padding
		n = pad(n, p)
	}

	if len(s.bracketBuf) > 0 {
		if (field == "disc" && s.multDiscs && n != "") || (field != "disc" && n != "") {
			if (*s).renderRight {
				n = n + string((*s).bracketBuf)
			} else {
				n = string((*s).bracketBuf) + n
			}
		}
	}
	Println(channel, "jdjdjdjdj\n"+n)
	Println(channel, n)
	return n
}

func pad(s string, a int) string {
	for len(s) > a {
		s = "0" + s
	}
	return s
}
