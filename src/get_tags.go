package main

import (
    "errors"
    "os"
    "path"
    "regexp" 
    "sort"
    "strconv"   
    "strings"
    "time"

    "github.com/Akumzy/ipc"
    "github.com/dhowden/tag"
    "github.com/gabriel-vasile/mimetype"
    ffmpeg "github.com/modfy/fluent-ffmpeg"
)

func getTags(channel *ipc.IPC, formData FormData, ffprobePath string) VideoData {
    ffmpeg.SetFfProbePath(ffprobePath)
    validatePaths(channel, formData)

    setLabel(channel, "reading " + path.Base(formData.albumDirectory) + "..")
    albumDirectoryFile, err := os.Open(formData.albumDirectory); if err != nil { panic(err) }
    files, err := albumDirectoryFile.Readdirnames(-1); if err != nil { panic(err) }
    
    discTracks := make(map[uint64]uint64)
    audioFiles := []AudioFile{}
    imageFiles := []string{}

    // wrote this myself :D i'll probably have to change it sometime
    trackRe := regexp.MustCompile(`^([0-9]+|[A-Za-z]{1,2}|[0-9]+[A-Za-z]|)(-| - |_| |)([0-9]+|[A-Za-z])(. | |_)`)

    for _, base := range files {
        setLabel(channel, "reading " + base + "..")
        file := path.Join(formData.albumDirectory, base)
        mime, err := mimetype.DetectFile(file); if err != nil { panic(err) }

        switch strings.Split(mime.String(), "/")[0] {
        case "audio":
            var disc, track uint64
            var artist, albumArtist, title, album string
            ffprobeData, err := ffmpeg.Probe(file); if err != nil { panic(err) }
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

            seconds, err := strconv.ParseFloat(ffprobeData["format"].(map[string]interface{})["duration"].(string), 64); if err != nil { panic(err) }

            // frankly i dont trust the tag library's assesment of track / disc numbers SORRY
            if (ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["disc"] != nil) {
                disc = parseTrackTag(ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["disc"].(string))
            } else { disc = 1 }

            if (ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["track"] != nil) {
                track = parseTrackTag(ffprobeData["format"].(map[string]interface{})["tags"].(map[string]interface{})["track"].(string))
            } else {
                if !trackRe.MatchString(file) {
                    panic(errors.New("please make sure your filenames start with a track number" +
                        "if they are not tagged properly (which would be preferrable). for exact" + 
                        "specifications as to what does and does not get detected as a track" + 
                        "number see https://github.com/sunglasseds/album2video"))
                }
                track, disc = parseTrack(file, trackRe)
            }
            
            if dt, ok := discTracks[disc]; ok {
                if dt < track {
                    discTracks[disc] = track                
                }
            } else {
                discTracks[disc] = track
            }

            audioFiles = append(audioFiles, AudioFile{
                filename: file,
                artist: artist,
                albumArtist: albumArtist,
                title: title,
                album: album,
                year: year,
                track: track,
                disc: disc,
                time: time.Duration(seconds * float64(time.Second)),
                cover: cover,
                discTracks: &discTracks,
            })

            case "image":
                imageFiles = append(imageFiles, file)
        }
    }    

    setLabel(channel, "ordering audio files..")
    sort.Sort(byTrack(audioFiles))

    Println(channel, audioFiles)

    return VideoData{
        formData: formData,
        audioFiles: audioFiles,
        imageFiles: imageFiles,
        discTracks: discTracks,
    }
}

func validatePaths(channel *ipc.IPC, formData FormData) {
    setLabel(channel, "validating album path..")
    stats, err := os.Stat(formData.albumDirectory); if err != nil { panic(err) }

    if stats.Mode().IsRegular() {
        formData.albumDirectory = path.Dir(formData.albumDirectory)
    } else if stats.IsDir() { /* pass */ } else {
        panic(errors.New("album directory is not a file or directory"))
    }

    if !formData.detectCover {
        setLabel(channel, "validating cover path..")
        stats, err = os.Stat(formData.coverPath); if err != nil { panic(err) }

        if stats.Mode().IsRegular() {
            formData.albumDirectory = path.Dir(formData.albumDirectory)
        } else {
            panic(errors.New(formData.coverPath + " is not a file"))
        }
    }

    setLabel(channel, "validating output path..")
    stats, err = os.Stat(formData.outputPath); if err != nil { panic(err) }

    if formData.separateFiles && stats.Mode().IsRegular() {
        formData.outputPath = path.Dir(formData.outputPath)
    } else if !formData.separateFiles && stats.IsDir() {
        formData.outputPath = path.Join(formData.outputPath, "out.mp4")
    } else if formData.separateFiles == stats.IsDir() { /* pass */ } else {
        panic(errors.New(formData.outputPath + " is not a file or directory"))
    }
}

func getMetadata(filename string) (tag.Metadata) {
    f, err := os.Open(filename); if err != nil { panic(err) }
    m, err := tag.ReadFrom(f); if err != nil { panic(err) }
    err = f.Close(); if err != nil { panic(err) }
    return m
}

func parseTrackTag(v string) (t uint64) {
    t, err := strconv.ParseUint(strings.Split(v, "/")[0], 10, 64)
    if err != nil { panic(err) }; return
}

func parseTrack(filename string, trackRe *regexp.Regexp) (d, t uint64) {
    submatches := trackRe.FindSubmatch([]byte(filename))
    discSubmatch := string(submatches[1])
    trackSubmatch := string(submatches[3])
    
    d, err := strconv.ParseUint(discSubmatch, 10, 64)
    if err != nil {
        d = 1
        switch len(discSubmatch) {
        case 1:
            d += uint64(strings.Index(strings.ToLower(string(discSubmatch[0])), "abcdefghijklmnopqrstuvwxyz"))
            fallthrough
        case 2:
            d += uint64(strings.Index(strings.ToLower(string(discSubmatch[1])), "abcdefghijklmnopqrstuvwxyz")*26)
        default:
            panic(err)
        }
    }

    t, err = strconv.ParseUint(trackSubmatch, 10, 64)
    if err != nil {
        if len(trackSubmatch) != 1 {
            panic(err)
        }
        t = uint64(strings.Index(strings.ToLower(string(trackSubmatch[0])), "abcdefghijklmnopqrstuvwxyz") + 1)
    }
    return
}

// le []AudioFile sorting interface
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

func overallTrackNumber(track, disc uint64, discTracks *map[uint64]uint64) (n uint64) {
    n = track
    for i := uint64(1); i <= disc-1; i++ {
        n += (*discTracks)[i]
    }
    return n
}
