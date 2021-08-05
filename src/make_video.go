package main

import (
	"bufio"
	"bytes"
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
	//"github.com/alessio/shellescape"
)

func makeVideo(channel *ipc.IPC, videoData VideoData, ffmpegPath string, ffprobePath string) {
	timestamps       := []Timestamp{}
	length           := time.Duration(0)
	fileListContents := ""
	
	for _, f := range videoData.audioFiles {
		timestamps = append(timestamps, Timestamp{
			artist: f.artist,
			title: f.title,
			time: durationToString(length),
		})
		
		fileListContents += "file '" + strings.ReplaceAll(f.filename, `'`, `'\''`) + "'\n"
		length += f.time
	}

	fileList, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), "tmp-*.txt")
	if err != nil { panic(err) }; defer os.Remove(fileList.Name())
	_, err = fileList.Write([]byte(fileListContents)); if err != nil { panic(err) }
	
	cw, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), "tmp-*.wav")
	if err != nil { panic(err) }; defer os.Remove(cw.Name())
	concatWavName := cw.Name()
	os.Remove(cw.Name()) // lol
	
	makeConcatWav := exec.Command(ffmpegPath, "-progress", "pipe:2", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	ffmpegStderr, _ := makeConcatWav.StderrPipe()
	makeConcatWav.Start()
	
	scanner := bufio.NewScanner(ffmpegStderr)
	scanner.Split(scanFFmpegChunks)
	for scanner.Scan() {
		m := scanner.Text()
		Println(channel, m)
	}
	
	makeConcatWav.Wait()
	
    //ffmpeg.SetFfProbePath(ffprobePath)
   	//stage1 := ffmpeg.NewCommand(ffmpegPath)
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
