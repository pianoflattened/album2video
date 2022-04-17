package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var re = regexp.MustCompile(`out_time_ms=(\d+)`)
var bad_size = regexp.MustCompile(`\[libx264 @ 0x[0-9a-f]+\] (height|width) not divisible by 2 \(([0-9]+)x([0-9]+)\)`)
const bufferSize = 65536

// lots of repeated code for running ffmpeg commands. please simplify
func makeVideo(bar ProgressBar, videoData VideoData, ffmpegPath string, fmtString string) string {
	timestamps := []Timestamp{}
	length := time.Duration(0)
	fileListContents := ""
	var c int

	// https://stackoverflow.com/questions/18434854/merge-m4a-files-in-terminal
	// all sequences of m4a files need to go through this process
	// then each concatenated sequence has to be converted to wav

	bar.Label = "making file list.."
	m4aSequence := []AudioFile{}
	
	for i, f := range videoData.audioFiles {
		//println(durationToString(length))
		if strings.HasSuffix(f.filename, ".m4a") {
			//fmt.Println(f.filename)
			m4aSequence = append(m4aSequence, f)
		} else {
			if len(m4aSequence) > 0 {
				concatMp3Name := concatM4aSequence(bar, m4aSequence, videoData)
				m4aSequence = []AudioFile{}
				// add concat mp3 to filelist
				fileListContents += "file '" + strings.ReplaceAll(concatMp3Name, `'`, `'\''`) + "'\n"
				// hope the defer doesnt trigger until the whole thing finishes :(
			}
			
			fileListContents += "file '" + strings.ReplaceAll(f.filename, `'`, `'\''`) + "'\n"
		}
		
		timestamps = append(timestamps, Timestamp{
			Artist:       f.artist,
			Title:        f.title,
			AlbumArtist:  f.albumArtist,
			Time:         durationToString(length),
			Disc:         f.disc,
			Track:        f.track,
			OverallTrack: i + 1,
		})

		length += f.time
	}
	
	if len(m4aSequence) > 0 {
		concatMp3Name := concatM4aSequence(bar, m4aSequence, videoData)
		fileListContents += "file '" + strings.ReplaceAll(concatMp3Name, `'`, `'\''`) + "'\n"
	}

	if videoData.formData.extractCover { // first try :D
		bar.Label = "extracting cover art.."
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

	bar.Label = "calculating minimum framerate.."
	framerate := float64((1.0 * float64(time.Second)) / float64(length))
	gop := framerate / 2.0
	//framerate = math.Max(framerate, 1.0)

	fileList, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[BIT_LY9099]--*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileList.Name())
	_, err = fileList.Write([]byte(fileListContents))
	//fmt.Println(fileListContents)
	if err != nil {
		panic(err)
	}

	bar = NewProgressBar(true)
	bar.Label = "concatenating audio files.."

	cw, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[BIT_LY9099]--*.wav")
	if err != nil {
		panic(err)
	}
	defer os.Remove(cw.Name())
	concatWavName := cw.Name()
	cw.Close()
	os.Remove(cw.Name()) // lol

	// ffmpeg -progress pipe:2 -y -f concat -safe 0 -i [filelist.txt] -c copy concat.wav
	//fmt.Println(ffmpegPath, "-progress", "pipe:2", "-y", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	makeConcatWav := exec.Command(ffmpegPath, "-progress", "pipe:2", "-y", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	concatStderr, _ := makeConcatWav.StderrPipe()
	makeConcatWav.Start()

	concatScanner := bufio.NewScanner(concatStderr)
	concatScanner.Split(scanFFmpegChunks)
	for concatScanner.Scan() {
		m := concatScanner.Text()
		a := re.FindAllStringSubmatch(m, -1)
		c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		//fmt.Println(c)
		bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
		fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
	}

	makeConcatWav.Wait()
	bar.Complete = true
	fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))

	fileList.Close()
	err = os.Remove(fileList.Name())
	if err != nil {
		fmt.Println(err)
	}

	bar = NewProgressBar(true)
	bar.Label = "making output video.."

	// ffmpeg -progress pipe:2 -y -loop 0 -r [framerate] -i [cover image] -i [concat.wav] -tune stillimage -t [length] -r [framerate] -c:a aac -profile:a aac_low -b:a 384k -pix_fmt yuv420p -c:v libx264 -profile:v high -preset slow -crf 18 -g [framerate/2] -movflags faststart [out.mp4]

	// normalOptions := []string{"-progress", "pipe:2", "-y", "-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", videoData.formData.coverPath, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate)}
	normalOptions := []string{"-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", videoData.formData.coverPath, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate)}

	youtubeOptions := []string{"-c:a", "aac", "-profile:a", "aac_low", "-b:a", "384k", "-pix_fmt", "yuv420p", "-c:v", "libx264", "-profile:v", "high", "-preset", "slow", "-crf", "18", "-g", fmt.Sprintf("%v", gop), "-movflags", "faststart"}

	coverCrop := doFFmpeg(bar, length, ffmpegPath, append(normalOptions, append(youtubeOptions, videoData.formData.outputPath)...)...)
	if len(coverCrop) > 0 {
		normalOptions := []string{"-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", coverCrop, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate)}
		doFFmpeg(bar, length, ffmpegPath, append(normalOptions, append(youtubeOptions, videoData.formData.outputPath)...)...)
		os.Remove(coverCrop)
	}

	if videoData.formData.extractCover {
		os.Remove(videoData.formData.coverPath)
	}
	
	fmt.Println(formatTimestamps(fmtString, timestamps))

	return concatWavName
}

func concatM4aSequence(bar ProgressBar, m4aSequence []AudioFile, videoData VideoData) (concatMp3Name string) {
	bar.Label = "m4a files found. extracting streams.."
	aacSequence := []string{}
				
	// extract all the aac streams
	concatM4aLength := time.Duration(0)
	for _, m4a := range m4aSequence {
		bar.Label = "extracting stream from '" + filepath.Base(m4a.filename) + "'.."
		aacSequence = append(aacSequence, strings.TrimSuffix(m4a.filename, "m4a")+"aac")
		
		doFFmpeg(bar, m4a.time, ffmpegPath, "-i", m4a.filename, "-acodec", "copy", strings.TrimSuffix(m4a.filename, "m4a")+"aac")
		concatM4aLength += m4a.time
		bar = NewProgressBar(true)
	}

	//aacSequenceStr := ""
	//for _, aac := range aacSequence {
	//	aacSequenceStr += shEscape(aac) + " "
	//}
	
	// concat them
	concatAacName := getConcatName(videoData, ".CONCAT--[GOO_GL6066]--*.aac")
	bar.Label = "concatenating aac files.."
	
	//var o []byte
	
	f, err := os.OpenFile(concatAacName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	for _, aac := range aacSequence {
		g, err := os.Open(aac)
		if err != nil {
			panic(err)
		}
		defer g.Close()
		
		buf := make([]byte, bufferSize)
		for {
			n, err := g.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}
			
			if _, err := f.Write(buf[:n]); err != nil {
				panic(err)
			}
		}
		g.Close()
	}
	f.Close()
	
	bar.Label = "cleaning up.."
	//for _, aac := range aacSequence {
		//os.Remove(aac)
	//}

	// convert to m4a
	concatM4aName := getConcatName(videoData, ".CONCAT--[T_CO7077]--*.m4a")
	
	doFFmpeg(bar, concatM4aLength, ffmpegPath, "-i", concatAacName, "-acodec", "copy", "-bsf:a", "aac_adtstoasc", concatM4aName)
	bar = NewProgressBar(true)

	// convert to mp3
	// may be able to skip step above. fine for now
	bar.Label = "converting concatenated file to mp3.."
	concatMp3Name = getConcatName(videoData, ".CONCAT--[ADF_LY8088]--*.mp3")

	// -q:a 0 is minimal quality loss. raise if it takes a bit too long but nothing past 4
	doFFmpeg(bar, concatM4aLength, ffmpegPath, "-i", concatM4aName, "-c:v", "copy", "-c:a", "libmp3lame", "-q:a", "0", concatMp3Name)
	bar = NewProgressBar(true)
	return
}

func doFFmpeg(bar ProgressBar, length time.Duration, ffmpegPath string, args ...string) string {
	ffmpegCmd := exec.Command(ffmpegPath, append([]string{"-progress", "pipe:2", "-y"}, args...)...)
	ffmpegStderr, _ := ffmpegCmd.StderrPipe()
	ffmpegCmd.Start()

	var c int
	ffmpegScanner := bufio.NewScanner(ffmpegStderr)
	ffmpegScanner.Split(scanFFmpegChunks)
	for ffmpegScanner.Scan() {
		m := ffmpegScanner.Text()
		if bad_size.MatchString(m) {
			dim := bad_size.FindAllStringSubmatch(m, 1)[0]
			fmt.Printf("%#v\n", dim)
			w, _ := strconv.Atoi(dim[2])
			h, _ := strconv.Atoi(dim[3])
			return fixDimensions(ffmpegPath, args[5], w, h) // bam
		}
		a := re.FindAllStringSubmatch(m, -1)
		c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
		fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
	}
	
	ffmpegCmd.Wait()
	bar.Complete = true
	fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))
	return ""
}

func getConcatName(videoData VideoData, format string) (concatName string) {
	c, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), format)
	if err != nil {
		panic(err)
	}
	defer os.Remove(c.Name())
	
	concatName = c.Name()
	c.Close()
	os.Remove(c.Name()) // lol
	return
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

// i am such a genius for writing this. wow. i am so cool
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

// gross
func shEscape(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "\\", "\\\\"), "\"", "\\\""), "'", "\\'"), " ", "\\ ")
}

func fixDimensions(ffmpegPath, imgPath string, w, h int) string {
	if h%2 == 1 { h -= 1 }
	if w%2 == 1 { w -= 1 }
	ext := filepath.Ext(imgPath)
	cropName := ".__" + strings.TrimSuffix(imgPath, ext) + "_crop" + ext
	
	ffmpegCmd := exec.Command(ffmpegPath, "-y", "-i", imgPath, "-vf", fmt.Sprintf("crop=%d:%d:0:0", w, h), cropName)
	fmt.Println([]string{ffmpegPath, "-y", "-i", imgPath, "-vf", fmt.Sprintf("\"crop=%d:%d:0:0\"", w, h), cropName})
	o, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(o))
		panic(err)
	}
	
	return cropName
}
