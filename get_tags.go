package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	//"github.com/Akumzy/ipc"
	"github.com/dhowden/tag"
	"github.com/gabriel-vasile/mimetype"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
)

// i was on crack and drugs when i wrote this
const alphabet = "abcdefghijklmnopqrstuvwxyz"

func getTags(bar ProgressBar, formData FormData, ffprobePath string) VideoData {
	ffmpeg.SetFfProbePath(ffprobePath)
	formData = validatePaths(bar, formData)

	bar.Label = "reading "+path.Base(formData.albumDirectory)+".."
	albumDirectoryFile, err := os.Open(formData.albumDirectory)
	if err != nil {
		panic(err)
	}
	files, err := albumDirectoryFile.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	discTracks := make(map[int]int)
	audioFiles := []AudioFile{}
	imageFiles := []string{}

	// wrote this myself :D i'll probably have to change it sometime
	trackRe := regexp.MustCompile(`^([0-9]+|[A-Za-z]{1,2}|[0-9]+[A-Za-z]|)([-_ ]| - |)([0-9]+|[A-Za-z])[ _.]`)

	for _, base := range files {
		if strings.HasPrefix(base, ".CONCAT--") {
			os.Remove(path.Join(formData.albumDirectory, base)) // clean up after yrself
			continue
		}
		bar.Label = "reading "+base+".."
		file := path.Join(formData.albumDirectory, base)
		mime, err := mimetype.DetectFile(file)
		if err != nil {
			if !(strings.HasSuffix(errors.Unwrap(err).Error(), "The handle is invalid.") || 
				 strings.HasSuffix(errors.Unwrap(err).Error(), "is a directory")) {
				panic(err)
			}
			continue
		}

		switch strings.Split(mime.String(), "/")[0] {
		case "audio":
			var disc, track int
			var artist, albumArtist, title, album string
			ffprobeData, err := ffmpeg.Probe(file)
			if err != nil {
				panic(err)
			}
			metadata := getMetadata(file)

			if artist = metadata.Artist(); artist == "" {
				artist = "[unknown artist]"
			}

			if albumArtist = metadata.AlbumArtist(); albumArtist == "" {
				albumArtist = "[unknown artist]"
			}

			if title = metadata.Title(); title == "" {
				title = "[untitled]"
			}

			if album = metadata.Album(); album == "" {
				album = "[untitled]"
			}

			year := strconv.Itoa(metadata.Year())
			cover := metadata.Picture()

			seconds_, err := strconv.ParseFloat(ffprobeData["format"].(map[string]interface{})["duration"].(string), 32)
			if err != nil {
				panic(err)
			}
			seconds := float64(seconds_)

			// frankly i dont trust the tag library's assesment of track / disc numbers SORRY
			if ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["disc"] != nil {
				disc = parseTrackTag(ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["disc"].(string))
			} else {
				disc = 1
			}

			if ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["track"] != nil {
				track = parseTrackTag(ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["track"].(string))
			} else {
				if !trackRe.MatchString(path.Base(file)) {
					panic(errors.New("please make sure your filenames start with a track number if they " +
						"are not tagged properly (which would be preferrable). for exact specifications as " +
						"to what does and does not get detected as a track number see https://github.com/" +
						"pianoflattened/album2video. if yr files are properly tagged and the track numbers " +
						"work then ffprobe is fucked and i have to use straight command output probably"))
				}
				track, disc = parseTrack(path.Base(file), trackRe)
			}

			if dt, ok := discTracks[disc]; ok {
				if dt < track {
					discTracks[disc] = track
				}
			} else {
				discTracks[disc] = track
			}

			audioFiles = append(audioFiles, AudioFile{
				filename:    file,
				artist:      artist,
				albumArtist: albumArtist,
				title:       title,
				album:       album,
				year:        year,
				track:       track,
				disc:        disc,
				time:        time.Duration(seconds * float64(time.Second)),
				cover:       cover,
				discTracks:  &discTracks,
			})

		case "image":
			imageFiles = append(imageFiles, file)
		}
	}

	bar.Label = "ordering audio files.."
	sort.Sort(byTrack(audioFiles))
	if len(audioFiles) < 1 {
		panic(errors.New("you need sound files in the album directory"))
	}

	return VideoData{
		formData:   formData,
		audioFiles: audioFiles,
		imageFiles: imageFiles,
		discTracks: discTracks,
	}
}

func validatePaths(bar ProgressBar, formData FormData) FormData {
	bar.Label = "validating album path.."
	stats, err := os.Stat(formData.albumDirectory)
	if err != nil {
		panic(err)
	}

	if stats.Mode().IsRegular() {
		formData.albumDirectory = path.Dir(formData.albumDirectory)
	} else if stats.IsDir() { /* pass */
	} else {
		panic(errors.New("album directory is not a file or directory"))
	}

	if !formData.extractCover {
		bar.Label = "validating cover path.."
		stats, err = os.Stat(formData.coverPath)
		if err != nil {
			panic(err)
		}

		if stats.Mode().IsRegular() { /* pass */
		} else {
			panic(errors.New(formData.coverPath + " is not a file"))
		}
	}

	bar.Label = "validating output path.."
	stats, err = os.Stat(formData.outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(formData.outputPath)
			if err != nil {
				panic(err)
			}
			f.Close()
		} else {
			panic(err)
		}
	} else {
		if formData.separateFiles && stats.Mode().IsRegular() {
			formData.outputPath = path.Dir(formData.outputPath)
		} else if (!formData.separateFiles) && stats.IsDir() {
			formData.outputPath = path.Join(formData.outputPath, "out.mp4")
		} else if formData.separateFiles == stats.IsDir() { /* pass */
		} else {
			panic(errors.New(formData.outputPath + " is not a file or directory"))
		}
	}
	return formData
}

func getMetadata(filename string) tag.Metadata {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	m, err := tag.ReadFrom(f)
	if err != nil {
		fmt.Println(filename)
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	return m
}

func parseTrackTag(v string) int {
	t, err := strconv.ParseInt(strings.Split(v, "/")[0], 10, 32)
	if err != nil {
		panic(err)
	}
	return int(t)
}

func parseTrack(filename string, trackRe *regexp.Regexp) (int, int) {
	submatches := trackRe.FindSubmatch([]byte(filename))
	discSubmatch := string(submatches[1])
	trackSubmatch := string(submatches[3])

	d, err := strconv.ParseInt(discSubmatch, 10, 32)
	if err != nil {
		d = 1
		switch len(discSubmatch) {
		case 2:
			d += int64(strings.Index(strings.ToLower(string(discSubmatch[1])), alphabet) * 26)
			fallthrough
		case 1:
			d += int64(strings.Index(strings.ToLower(string(discSubmatch[0])), alphabet))
		default:
			panic(err)
		}
	}

	t, err := strconv.ParseInt(trackSubmatch, 10, 32)
	if err != nil {
		if len(trackSubmatch) != 1 {
			panic(err)
		}
		t = int64(strings.Index(strings.ToLower(string(trackSubmatch[0])), alphabet) + 1)
	}

	return int(d), int(t)
}

// le []AudioFile sorting interface
// see line 19 for more info
type byTrack []AudioFile

func (s byTrack) Len() int {
	return len(s)
}

func (s byTrack) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byTrack) Less(i, j int) bool {
	discTracks := s[i].discTracks
	iOverall := overallTrackNumber(s[i].track, s[i].disc, discTracks)
	jOverall := overallTrackNumber(s[j].track, s[j].disc, discTracks)
	return iOverall < jOverall
}

func overallTrackNumber(track, disc int, discTracks *map[int]int) (n int) {
	n = track
	for i := 1; i <= disc-1; i++ {
		n += (*discTracks)[i]
	}
	return n
}
