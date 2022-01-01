package main

import (
	"fmt"
	"os"
	"math"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

var spinner = [6]string{"⠋", "⠙", "⠸", "⠴", "⠦", "⠇"}

type ProgressBar struct {
	Label    	string
	Progress 	float64
	Previous	int
	Determinate bool
	Complete	bool
}

func NewProgressBar(determinate bool) (bar ProgressBar) {
	return ProgressBar{
		Label: "",
		Progress: 0.0,
		Previous: -1,
		Determinate: determinate,
		Complete: false,
	}
}

func (p *ProgressBar) Render(currentSize, completeSize float64) (bar string) {
	bar = ""
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	
	if err != nil {
		width = 0
	}

	if !p.Determinate {
		minout := utf8.RuneCountInString(p.Label) + 2
		(*p).Previous = (p.Previous + 1) % 6
		spinnerStage := spinner[p.Previous]
		
		if p.Complete {
			spinnerStage = "✓"
		}
		
		if width - minout < 0 {
			if width < 1 {
				bar = ""
			}
			return spinnerStage
		}
		return fmt.Sprintf("%s %s\r", spinnerStage, p.Label)
		
	} else {
		minout := utf8.RuneCountInString(fmt.Sprintf("|░| 100%% (%.1f/%.1f MB)", completeSize, completeSize))
		pbarWidth := width-minout+1
		pbarAmount := int(math.Floor(float64(pbarWidth) * (*p).Progress))
		percentStr := fmt.Sprintf("%3.0f%%", (*p).Progress * 100)
		
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
						return percentStr
					}
					return fmt.Sprintf("%.1f/%.1f MB\r", currentSize, completeSize)
				}
				return fmt.Sprintf("(%.1f/%.1f MB)\r", currentSize, completeSize)
			}
			return fmt.Sprintf("%s (%.1f/%.1f MB)\r", percentStr, currentSize, completeSize)
		}
		return fmt.Sprintf("|%s%s| %s (%.1f/%.1f MB)\r", strings.Repeat("█", pbarAmount), strings.Repeat("░", pbarWidth-pbarAmount), percentStr, currentSize, completeSize)
	}
	return
}
