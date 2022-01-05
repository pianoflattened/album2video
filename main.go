package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pborman/getopt/v2"
)

var (
	albumDirectory = "."
	coverPath = ""
	extractCover bool
	separateFiles bool
	outputPath = "./out.mp4"
	verbose bool
	quiet bool
	ffprobePath = ""
	ffmpegPath = ""
	fmtString = "%[ - >a%s - %t"
	help bool
)

func init() {
	getopt.FlagLong(&verbose, "verbose", 'v', "make script say more things")
	getopt.FlagLong(&quiet, "quiet", 'q', "shut up !")
	getopt.FlagLong(&coverPath, "cover", 'c', "specify cover art. video will be black if unset")
	getopt.FlagLong(&extractCover, "extract-cover", 'x', "use cover art embedded in id3 tags")
	getopt.FlagLong(&separateFiles, "separate", 's', "create separate files, 1 per sound file")
	getopt.FlagLong(&albumDirectory, "album-directory", 'i', "directory containing the sound files. this is the current directory if not specified")
	getopt.FlagLong(&outputPath, "output-path", 'o', "the filename of the resulting video. if -s / --separate is used, this should be the directory all the video files will go in. ./out.mp4 / the current directory if not specified")
	getopt.FlagLong(&ffprobePath, "ffprobe", 0, "path to specific ffprobe binary")
	getopt.FlagLong(&ffmpegPath, "ffmpeg", 0, "path to specific ffmpeg binary")
	getopt.FlagLong(&fmtString, "track-format", 'f', "track formatting")
	getopt.FlagLong(&help, "help", 'h', "show help data")
}

func main() {
	getopt.Parse()
	
	if help {
		// !!!!!!!!!!! WHY WAS THIS SO HARD 2 FIGURE OUT
		getopt.Usage()
		return
	}
	
	formData := FormData{
		albumDirectory: albumDirectory,
		coverPath:      coverPath,
		extractCover:   extractCover,
		separateFiles:  separateFiles,
		outputPath:     outputPath,
		verbose:		verbose,
		quiet:			quiet,
	}
	
	if ffprobePath == "" {
		var b = []byte{}
		var err error
		
		switch os := runtime.GOOS; os {
		case "windows":
			b, err = exec.Command("where", "ffprobe").Output()
		default:
			b, err = exec.Command("which", "ffprobe").Output()
		}
		ffprobePath = strings.TrimSpace(string(b))
		
		if err != nil {
			fmt.Println("couldnt find ffprobe executable on your PATH. please specify with the --ffprobe option")
			panic(err)
		}
	}
	
	if ffmpegPath == "" {
		var b = []byte{}
		var err error
		
		switch os := runtime.GOOS; os {
		case "windows":
			b, err = exec.Command("where", "ffmpeg").Output()
		default:
			b, err = exec.Command("which", "ffmpeg").Output()
		}
		ffmpegPath = strings.TrimSpace(string(b))
		
		if err != nil {
			fmt.Println("couldnt find ffmpeg executable on your PATH. please specify with the --ffmpeg option")
			panic(err)
		}
	}
	
	bar := NewProgressBar(false)

	videoData := getTags(bar, formData, ffprobePath)
	/*concatWav*/_ = makeVideo(bar, videoData, ffmpegPath, fmtString)
}
