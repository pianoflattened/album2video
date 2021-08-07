package main

import (
	"math"
	"strconv"
	"strings"
)

type fmtState struct {
	tokenBuf   []rune
	padding    int
	fCase      FieldCase
	bracketBuf []rune
	renderRight,
	bracketEscaped,
	multDiscs,
	inBrackets bool
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

func (s *fmtState) render(args ...interface{}) string { // evil definition sorry ik its gross
	if len(args) > 2 {
		return ""
	}

	var field string
	time := ""

	if len(args) == 2 {
		if args[1].(bool) {
			field = args[0].(string)
			time = pad(strings.Join(strings.Split(field, ":"), ""), int(math.Max(float64((*s).padding), 3)))
			if len(time) > 4 {
				time = time[0:len(time)-4] + ":" + time[len(time)-4:]
			}
			s.clear()
			field = time[0:len(time)-2] + ":" + time[len(time)-2:]
		}
	}

	switch args[0].(type) {
	case string:
		field = args[0].(string)
		if (*s).fCase == lower {
			s.clear()
			field = strings.ToLower(field)
		} else if (*s).fCase == title {
			s.clear()
			field = TitleCase(field, false)
		}

	case int:
		p := (*s).padding
		s.clear()
		field = pad(strconv.Itoa(args[0].(int)), p)
	}

	if len(s.bracketBuf) > 0 {
		if field != "" {
			if len(args) > 1 {
				if !args[1].(bool) && s.multDiscs {
					goto renderside
				}
			}
		renderside:
			if s.renderRight {
				field = field + string(s.bracketBuf)
			} else {
				field = string(s.bracketBuf) + field
			}
		}
	}

	s.clear()

	return field
}

func pad(s string, a int) string {
	for len(s) > a {
		s = "0" + s
	}
	return s
}
