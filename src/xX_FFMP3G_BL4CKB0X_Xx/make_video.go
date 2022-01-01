package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func makeVideo(bar ProgressBar, videoData VideoData, ffmpegPath string, fmtString string) string {
	timestamps := []Timestamp{}
	length := time.Duration(0)
	fileListContents := ""

	setLabel("making file list..")
	for i, f := range videoData.audioFiles {
		println(durationToString(length))
		timestamps = append(timestamps, Timestamp{
			Artist:       f.artist,
			Title:        f.title,
			AlbumArtist:  f.albumArtist,
			Time:         durationToString(length),
			Disc:         f.disc,
			Track:        f.track,
			OverallTrack: i + 1,
		})

		fileListContents += "file '" + strings.ReplaceAll(f.filename, `'`, `'\''`) + "'\n"
		length += f.time
	}

	if videoData.formData.extractCover { // first try :D
		setLabel("extracting cover art..")
		if videoData.audioFiles[0].cover == nil {
			panic(errors.New("there is no cover art embedded into the first track. please tag your files properly"))
		}
		autoCoverFile, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".tmp-*."+videoData.audioFiles[0].cover.Ext)
		if err != nil {
			panic(err)
		}
		_, err = autoCoverFile.Write(videoData.audioFiles[0].cover.Data)
		if err != nil {
			panic(err)
		}
		defer os.Remove(autoCoverFile.Name())

		videoData.formData.coverPath = autoCoverFile.Name()
	}

	setLabel("calculating minimum framerate..")
	framerate := float64((1.0 * float64(time.Second)) / float64(length))
	gop := framerate / 2.0
	//framerate = math.Max(framerate, 1.0)

	fileList, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[BIT_LY9099]--*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileList.Name())
	_, err = fileList.Write([]byte(fileListContents))
	if err != nil {
		panic(err)
	}

	makeDeterminate()
	setLabel("concatenating audio files..")

	cw, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[BIT_LY9099]--*.wav")
	if err != nil {
		panic(err)
	}
	defer os.Remove(cw.Name())
	concatWavName := cw.Name()
	cw.Close()
	os.Remove(cw.Name()) // lol

	// ffmpeg -progress pipe:2 -y -f concat -safe 0 -i [filelist.txt] -c copy concat.wav
	makeConcatWav := exec.Command(ffmpegPath, "-progress", "pipe:2", "-y", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	concatStderr, _ := makeConcatWav.StderrPipe()
	makeConcatWav.Start()

	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	concatScanner := bufio.NewScanner(concatStderr)
	concatScanner.Split(scanFFmpegChunks)
	for concatScanner.Scan() {
		m := concatScanner.Text()
		a := re.FindAllStringSubmatch(m, -1)
		c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		setProgress(float64(time.Duration(c)*time.Microsecond)/float64(length))
	}

	makeConcatWav.Wait()

	fileList.Close()
	err = os.Remove(fileList.Name())
	if err != nil {
		Println(err)
	}

	makeDeterminate()
	setLabel("making output video..")

	// ffmpeg -progress pipe:2 -y -loop 0 -r [framerate] -i [cover image] -i [concat.wav] -tune stillimage -t [length] -r [framerate] -c:a aac -profile:a aac_low -b:a 384k -pix_fmt yuv420p -c:v libx264 -profile:v high -preset slow -crf 18 -g [framerate/2] -movflags faststart [out.mp4]

	normalOptions := []string{"-progress", "pipe:2", "-y", "-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", videoData.formData.coverPath, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate)}

	youtubeOptions := []string{"-c:a", "aac", "-profile:a", "aac_low", "-b:a", "384k", "-pix_fmt", "yuv420p", "-c:v", "libx264", "-profile:v", "high", "-preset", "slow", "-crf", "18", "-g", fmt.Sprintf("%v", gop), "-movflags", "faststart"}

	makeOutputVideo := exec.Command(ffmpegPath, append(normalOptions, append(youtubeOptions, videoData.formData.outputPath)...)...)
	outputStderr, _ := makeOutputVideo.StderrPipe()
	makeOutputVideo.Start()

	outputScanner := bufio.NewScanner(outputStderr)
	outputScanner.Split(scanFFmpegChunks)
	for outputScanner.Scan() {
		m := outputScanner.Text()
		//println(m)
		a := re.FindAllStringSubmatch(m, -1)
		c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		setProgress(float64(time.Duration(c)*time.Microsecond)/float64(length))
	}

	makeOutputVideo.Wait()

	if videoData.formData.extractCover {
		os.Remove(videoData.formData.coverPath)
	}

	formatTimestamps(timestamps, fmtString)
	setComplete()

	return concatWavName
}

func durationToString(d time.Duration) (t string) {
	timeSlice := regexp.MustCompile(`[hm]`).Split(strings.TrimRight(d.String(), "s"), 3)
	var hours, minutes int
	var seconds float64
	switch len(timeSlice) {
	case 1:
		seconds, _ = strconv.ParseFloat(timeSlice[0], 32)
	case 2:
		minutes, _ = strconv.Atoi(timeSlice[0])
		seconds, _ = strconv.ParseFloat(timeSlice[1], 32)
	default:
		hours, _ = strconv.Atoi(timeSlice[0])
		minutes, _ = strconv.Atoi(timeSlice[1])
		seconds, _ = strconv.ParseFloat(timeSlice[2], 32)
	}
	t = strings.Split(fmt.Sprintf("%02d:%02d:%02d", hours, minutes, int(math.Floor(seconds))), ".")[0]
	for (strings.HasPrefix(t, "0") || strings.HasPrefix(t, ":")) && len(t) > 4 {
		t = trimLeftChar(t)
	}
	return
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func scanFFmpegChunks(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte("help")); i >= 0 {
		return i + 4, nil, nil
	}

	if i := bytes.Index(data, []byte("progress=continue")); i >= 0 {
		return i + 17, data[0:i], nil
	}

	if i := bytes.Index(data, []byte("progress=end")); i >= 0 {
		return i + 12, data[0:i], bufio.ErrFinalToken
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
