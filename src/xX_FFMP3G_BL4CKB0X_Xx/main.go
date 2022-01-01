package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/pborman/getopt/v2"
)

var {
	albumDirectory = "."
	coverPath = ""
	extractCover bool
	separateFiles bool
	outputPath = "./out.mp4"
	verbose bool
	quiet bool
	ffprobePath = ""
	ffmpegPath = ""
	fmtString = ""
	help bool

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
	getopt.FlagLong(&help, "h", "--help", "show help data")
}

func main() {
	if help {
		// !!!!!!!!!!! WHY IS THIS SO HARD 2 FIGURE OUT
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
		var err error
		
		switch os := runtime.GOOS; os {
		case "windows":
			ffprobePath, err = exec.Command("where ffprobe").Output()
		default:
			ffprobePath, err = exec.Command("which ffprobe").Output()
		}
		
		if err != nil {
			fmt.Println("couldnt find ffprobe executable on your PATH. please specify with the --ffprobe option")
			panic(err)
		}
	}
	
	if ffmpegPath == "" {
		var err error
		
		switch os := runtime.GOOS; os {
		case "windows":
			ffmpegPath, err = exec.Command("where ffmpeg").Output()
		default:
			ffmpegPath, err = exec.Command("which ffmpeg").Output()
		}
		
		if err != nil {
			fmt.Println("couldnt find ffmpeg executable on your PATH. please specify with the --ffmpeg option")
			panic(err)
		}
	}

	videoData := getTags(formData, ffprobePath)
	/*concatWav*/_ = makeVideo(videoData, ffmpegPath, fmtString)
}
