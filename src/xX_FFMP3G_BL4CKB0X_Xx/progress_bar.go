package main

import (
	"fmt"
	"os"
	"math"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const spinner = [6]string["⠋", "⠙", "⠸", "⠴", "⠦", "⠇"]

type ProgressBar struct {
	Label    	string
	Progress 	float64
	Previous	float64
	Determinate bool
	Complete	bool
}

func (p *ProgressBar) Render(currentSize, completeSize float64) (bar string) {
	bar = ""
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	
	if err != nil {
		width = 0
	}

	if !p.Determinate {
		braille := 
		if width < 1 {
			bar = ""
		}
		
		if p.Complete {
			bar = "✓"
		} else {
			bar = spinner[p.Previous + 1]
			(*p).Previous += 1
		}
		
	} else {
		minout := utf8.RuneCountInString(fmt.Sprintf("|░| 100%% (%.1f/%.1f MB)", completeSize, completeSize))
		pbarWidth := width-minout+1
		percentComplete := currentSize / completeSize
		pbarAmount := math.Floor(pbarWidth * percentComplete)
		percentStr := fmt.Sprintf("%3.0f%%", percentComplete * 100)
		
		if p.Complete {
			currentSize = completeSize
			percentStr = "100%"
			pbarAmount = pbarWidth
		}
		
		if width < minout {
			if width < minout - 4 { // len(fmt.Sprintf("100%% (%.1f/%.1f MB)", completeSize, completeSize))
				if width < minout - 9 { // len(fmt.Sprintf("(%.1f/%.1f MB)", completeSize, completeSize))
					if width < minout - 11 { // len(fmt.Sprintf("%.1f/%.1f MB", completeSize, completeSize))
						if width < 4 { // len("100%")
							bar = ""
						}
						bar = percentStr
					}
					bar = fmt.Sprintf("%.1f/%.1f MB", currentSize, completeSize)
				}
				bar = fmt.Sprintf("(%.1f/%.1f MB)", currentSize, completeSize)
			}
			bar = fmt.Sprintf("%s (%.1f/%.1f MB)", percentStr, completeSize, completeSize)
		}
		bar = fmt.Sprintf("|%s| %s (%.1f/%.1f MB)", strings.Repeat("█", pbarAmount), percentStr, completeSize, completeSize)
	}
	return
}
