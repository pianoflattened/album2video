package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path"
    "regexp"
    "strconv"
    "strings"
    "time"

	"github.com/Akumzy/ipc"
)

func makeVideo(channel *ipc.IPC, videoData VideoData, ffmpegPath string, ffprobePath string) {
	timestamps       := []Timestamp{}
	length           := time.Duration(0)
	fileListContents := ""
	
	setLabel(channel, "making file list..")
	for _, f := range videoData.audioFiles {
		timestamps = append(timestamps, Timestamp{
			artist: f.artist,
			title: f.title,
			time: durationToString(length),
		})
		
		fileListContents += "file '" + strings.ReplaceAll(f.filename, `'`, `'\''`) + "'\n"
		length += f.time
	}
	
	if videoData.formData.detectCover {
		setLabel(channel, "extracting cover art..")
		if videoData.audioFiles[0].cover == nil {
			panic(errors.New("there is no cover art embedded into the first track. please tag your files properly"))
		}
		autoCoverFile, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".tmp-*."+videoData.audioFiles[0].cover.Ext)
		if err != nil { panic(err) }
		_, err = autoCoverFile.Write(videoData.audioFiles[0].cover.Data); if err != nil { panic(err) }
		defer os.Remove(autoCoverFile.Name())
		
		videoData.formData.coverPath = autoCoverFile.Name()
	}
	
	setLabel(channel, "calculating minimum framerate..")
	framerate := float32((1.0*float32(time.Second))/float32(length))

	fileList, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".tmp-*.txt")
	if err != nil { panic(err) }; defer os.Remove(fileList.Name())
	_, err = fileList.Write([]byte(fileListContents)); if err != nil { panic(err) }
	
	makeDeterminate(channel)
	setLabel(channel, "concatenating audio files..")
	
	cw, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".tmp-*.wav")
	if err != nil { panic(err) }; defer os.Remove(cw.Name())
	concatWavName := cw.Name()
	os.Remove(cw.Name()) // lol
	
	makeConcatWav := exec.Command(ffmpegPath, "-progress", "pipe:2", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	concatStderr, _ := makeConcatWav.StderrPipe()
	makeConcatWav.Start()
	
	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	concatScanner := bufio.NewScanner(concatStderr)
	concatScanner.Split(scanFFmpegChunks)
	for concatScanner.Scan() {
		m := concatScanner.Text()
		a := re.FindAllStringSubmatch(m, -1)
		c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		setProgress(channel, float32(time.Duration(c)*time.Microsecond) / float32(length))
	}
	
	makeConcatWav.Wait()
	
	makeDeterminate(channel)
	setLabel(channel, "making output video..")
	makeOutputVideo := exec.Command(ffmpegPath, "-progress", "pipe:2", "-y", "-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", videoData.formData.coverPath, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate), "-c", "copy", videoData.formData.outputPath)
	outputStderr, _ := makeOutputVideo.StderrPipe()
	makeOutputVideo.Start()
	
	outputScanner := bufio.NewScanner(outputStderr)
	outputScanner.Split(scanFFmpegChunks)
	for outputScanner.Scan() {
		m := outputScanner.Text()
		a := re.FindAllStringSubmatch(m, -1)
		c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		setProgress(channel, float32(time.Duration(c)*time.Microsecond) / float32(length))
	}
	
	makeOutputVideo.Wait()

   	if videoData.formData.detectCover {
   		os.Remove(videoData.formData.coverPath)
   	}
}

func durationToString(d time.Duration) (t string) {
	timeSlice := regexp.MustCompile(`[hm]`).Split(strings.TrimRight(d.String(), "s"), 3)
	var hours, minutes, seconds int
	switch len(timeSlice) {
	case 1:
		seconds, _ = strconv.Atoi(timeSlice[0])
	case 2:
		minutes, _ = strconv.Atoi(timeSlice[0])
		seconds, _ = strconv.Atoi(timeSlice[1])
	default:
		hours, _   = strconv.Atoi(timeSlice[0])
		minutes, _ = strconv.Atoi(timeSlice[1])
		seconds, _ = strconv.Atoi(timeSlice[2])	
	}
	t = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
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
