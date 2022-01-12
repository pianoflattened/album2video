package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// lots of repeated code for running ffmpeg commands. please simplify
func makeVideo(bar ProgressBar, videoData VideoData, ffmpegPath string, fmtString string) string {
	timestamps := []Timestamp{}
	length := time.Duration(0)
	fileListContents := ""
	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	var c int

	// https://stackoverflow.com/questions/18434854/merge-m4a-files-in-terminal
	// all sequences of m4a files need to go through this process
	// then each concatenated sequence has to be converted to wav

	bar.Label = "making file list.."
	m4aSequence := []string{}
	aacSequence := []string{}
	
	var extractAacStream *exec.Cmd
	var aacStderr io.ReadCloser
	var aacScanner *bufio.Scanner
	
	var aacToM4a *exec.Cmd
	var aacToM4aStderr io.ReadCloser
	var aacToM4aScanner *bufio.Scanner

	var m4aToMp3 *exec.Cmd
	var m4aToMp3Stderr io.ReadCloser
	var m4aToMp3Scanner *bufio.Scanner
	
	for i, f := range videoData.audioFiles {
		//println(durationToString(length))
		fmt.Println(f.filename)
		if strings.HasSuffix(f.filename, ".m4a") {
			m4aSequence = append(m4aSequence, f.filename)
		} else {
			if len(m4aSequence) > 0 {
				bar.Label = "m4a files found. extracting streams.."
				
				// extract all the aac streams
				for _, m4a := range m4aSequence {
					bar.Label = "extracting stream from '" + filepath.Base(m4a) + "'.."
					aacSequence = append(aacSequence, strings.TrimSuffix(m4a, "m4a")+"aac")
					extractAacStream = exec.Command(ffmpegPath, "-progress", "pipe:2", "-y", "-i", m4a, "-acodec", "copy", strings.TrimSuffix(m4a, "m4a")+"aac")
					aacStderr, _ = extractAacStream.StderrPipe()
					extractAacStream.Start()

					aacScanner = bufio.NewScanner(aacStderr)
					aacScanner.Split(scanFFmpegChunks)
					for aacScanner.Scan() {
						m := aacScanner.Text()
						a := re.FindAllStringSubmatch(m, -1)
						c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
						bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
						fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
					}
					
					extractAacStream.Wait()
					bar.Complete = true
					fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))
					
					bar = NewProgressBar(true)
				}
				
				m4aSequence = []string{}

				aacSequenceStr := ""
				for _, aac := range aacSequence {
					aacSequenceStr += shEscape(aac) + " "
				}
				
				// concat them
				ca, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[GOO_GL6066]--*.aac")
				if err != nil {
					panic(err)
				}
				defer os.Remove(ca.Name())
				concatAacName := ca.Name()
				ca.Close()
				os.Remove(ca.Name()) // lol

				bar.Label = "concatenating aac files.."
				err = exec.Command("sh", "-c", "cat "+aacSequenceStr+"> "+concatAacName).Run()
				if err != nil {
					panic(err)
				}

				// convert to m4a
				cm, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[T_CO7077]--*.m4a")
				if err != nil {
					panic(err)
				}
				defer os.Remove(cm.Name())
				concatM4aName := cm.Name()
				cm.Close()
				os.Remove(cm.Name()) // lol
				
				aacToM4a = exec.Command(ffmpegPath, "-progress", "-pipe:2", "-y", "-i", concatAacName, "-acodec", "copy", "-bsf:a", "aac_adtstoasc", concatM4aName)
				aacToM4aStderr, _ = aacToM4a.StderrPipe()
				aacToM4a.Start()

				aacToM4aScanner = bufio.NewScanner(aacToM4aStderr)
				aacToM4aScanner.Split(scanFFmpegChunks)
				for aacToM4aScanner.Scan() {
					m := aacToM4aScanner.Text()
					a := re.FindAllStringSubmatch(m, -1)
					c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
					bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
					fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
				}

				aacToM4a.Wait()
				bar.Complete = true
				fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))

				// convert to mp3
				// may be able to skip step above. fine for now
				bar.Label = "converting concatenated file to mp3.."
				cmm, err := ioutil.TempFile(path.Dir(videoData.audioFiles[0].filename), ".CONCAT--[ADF_LY8088]--*.mp3")
				if err != nil {
					panic(err)
				}
				defer os.Remove(cmm.Name())
				concatMp3Name := cmm.Name()
				cmm.Close()
				os.Remove(cmm.Name()) // lol

				// -q:a 0 is minimal quality loss. raise if it takes a bit too long but nothing past 4
				m4aToMp3 = exec.Command(ffmpegPath, "-progress", "-pipe:2", "-y", "-i", concatM4aName, "-c:v", "copy", "-c:a", "libmp3lame", "-q:a", "0", concatMp3Name)
				m4aToMp3Stderr, _ = m4aToMp3.StderrPipe()
				m4aToMp3.Start()

				m4aToMp3Scanner = bufio.NewScanner(m4aToMp3Stderr)
				m4aToMp3Scanner.Split(scanFFmpegChunks)
				for m4aToMp3Scanner.Scan() {
					m := m4aToMp3Scanner.Text()
					a := re.FindAllStringSubmatch(m, -1)
					c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
					bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
					fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
				}

				m4aToMp3.Wait()
				bar.Complete = true
				fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))

				// add concat mp3 to filelist
				fileListContents += "file '" + strings.ReplaceAll(concatMp3Name, `'`, `'\''`) + "'\n"
				// hope the defer doesnt trigger until the whole thing finishes :(
			}
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

		fileListContents += "file '" + strings.ReplaceAll(f.filename, `'`, `'\''`) + "'\n"
		length += f.time
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
	makeConcatWav := exec.Command(ffmpegPath, "-progress", "pipe:2", "-y", "-f", "concat", "-safe", "0", "-i", fileList.Name(), "-c", "copy", concatWavName)
	concatStderr, _ := makeConcatWav.StderrPipe()
	makeConcatWav.Start()

	concatScanner := bufio.NewScanner(concatStderr)
	concatScanner.Split(scanFFmpegChunks)
	for concatScanner.Scan() {
		m := concatScanner.Text()
		a := re.FindAllStringSubmatch(m, -1)
		c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
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

	normalOptions := []string{"-progress", "pipe:2", "-y", "-loop", "0", "-r", fmt.Sprintf("%v", framerate), "-i", videoData.formData.coverPath, "-i", concatWavName, "-t", fmt.Sprintf("%v", length.Seconds()), "-r", fmt.Sprintf("%v", framerate)}

	youtubeOptions := []string{"-c:a", "aac", "-profile:a", "aac_low", "-b:a", "384k", "-pix_fmt", "yuv420p", "-c:v", "libx264", "-profile:v", "high", "-preset", "slow", "-crf", "18", "-g", fmt.Sprintf("%v", gop), "-movflags", "faststart"}

	makeOutputVideo := exec.Command(ffmpegPath, append(normalOptions, append(youtubeOptions, videoData.formData.outputPath)...)...)
	outputStderr, _ := makeOutputVideo.StderrPipe()
	makeOutputVideo.Start()

	outputScanner := bufio.NewScanner(outputStderr)
	outputScanner.Split(scanFFmpegChunks)

	for outputScanner.Scan() {
		m := outputScanner.Text()
		//fmt.Println(m)
		a := re.FindAllStringSubmatch(m, -1)
		c, _ = strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
		bar.Progress = float64(time.Duration(c)*time.Microsecond)/float64(length)
		fmt.Print((&bar).Render(time.Duration(c)*time.Microsecond, length))
	}

	makeOutputVideo.Wait()

	if videoData.formData.extractCover {
		os.Remove(videoData.formData.coverPath)
	}

	bar.Complete = true
	fmt.Println((&bar).Render(time.Duration(c)*time.Microsecond, length))
	fmt.Println(formatTimestamps(fmtString, timestamps))

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
