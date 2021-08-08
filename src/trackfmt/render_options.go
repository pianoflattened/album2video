package main

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

type renderOptions struct {
	fCase          FieldCase
	padding        int
	ifExistsRight  bool
	ifExists       string
	mode           rune
	dominantArtist string
	multDiscs      bool
}

func (r renderOptions) render(t Timestamp) string {
	s := r.valFromMode(t)

	if strings.IndexRune("tdnw", r.mode) > -1 { // padding
		if r.mode == 't' {
			s = pad(strings.ReplaceAll(s, ":", ""), int(math.Max(float64(r.padding), 3)))
			if len(s) > 4 {
				s = s[0:len(s)-4] + ":" + s[len(s)-4:]
			}
			s = s[0:len(s)-2] + ":" + s[len(s)-2:]
		} else {
			s = pad(s, r.padding)
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
		s = TitleCase(s)
	}

	return s
}

// time 2 whip out th reflect library for opaque bullshit >:)
func (r renderOptions) valFromMode(t Timestamp) string {
	switch f := reflect.ValueOf(t).FieldByName(string([]rune{r.mode})); f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(f.Int(), 10)
	case reflect.String:
		return f.String()
	}
	return ""
}

func pad(s string, a int) string {
	for len(s) < a {
		s = "0" + s
	}
	return s
}
