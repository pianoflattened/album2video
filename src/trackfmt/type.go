package main

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
)

type JSONTimestamp struct {
	Artist       string `json:"artist"`
	AlbumArtist  string `json:"albumArtist"`
	Title        string `json:"title"`
	Time         string `json:"time"`
	Disc         int    `json:"disc"`
	Track        int    `json:"track"`
	OverallTrack int    `json:"overallTrack"`
}

func (j JSONTimestamp) toTimestamp() Timestamp {
	return Timestamp{j.Artist, j.Artist, j.AlbumArtist, j.Title, j.Time, j.Disc, j.Track, j.OverallTrack}
}

type Timestamp struct {
	r           string
	a           string
	albumArtist string
	s           string
	t           string
	d           int
	n           int
	w           int
}

type renderOptions struct {
	fCase         FieldCase
	padding       int
	ifExistsRight bool
	ifExists      string
	mode          rune
}

func (r renderOptions) render(t Timestamp, ifExists string) string {
	s := r.valFromMode(t)

	if r.fCase == lower { // case
		s = strings.ToLower(s)
	} else if r.fCase == title {
		s = TitleCase(s, false)
	}

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

	if r.ifExists != "" {

	}

	return s
}

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
	for len(s) > a {
		s = "0" + s
	}
	return s
}
