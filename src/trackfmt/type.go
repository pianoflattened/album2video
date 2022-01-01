package main

type FieldCase int

const (
	raw FieldCase = iota
	lower
	title
	upper
)

func (t Timestamp) toFmtTimestamp() FmtTimestamp {
	return Timestamp{t.Artist, t.Artist, t.Title, t.Time, t.Disc, t.Track, t.OverallTrack}
}

type FmtTimestamp struct {
	r, a, s, t string
	d, n, w    int
}
