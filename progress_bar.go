package main

import (
	"fmt"
	"os"
	"math"
	"strings"
	"time"
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

func (p *ProgressBar) Render(_currentSize, _completeSize time.Duration) (bar string) {
	// grr
	currentSize := fmtDuration(_currentSize)
	completeSize := fmtDuration(_completeSize)
	
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
		minout := utf8.RuneCountInString(fmt.Sprintf("|░| 100%% (%s/%s)", completeSize, completeSize))
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
					return fmt.Sprintf("%s/%s\r", currentSize, completeSize)
				}
				return fmt.Sprintf("(%s/%s)\r", currentSize, completeSize)
			}
			return fmt.Sprintf("%s (%s/%s)\r", percentStr, currentSize, completeSize)
		}
		return fmt.Sprintf("|%s%s| %s (%s/%s)\r", strings.Repeat("█", pbarAmount), strings.Repeat("░", pbarWidth-pbarAmount), percentStr, currentSize, completeSize)
	}
	return
}

func fmtDuration(d time.Duration) string {
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    
    if h < 1 {
		return fmt.Sprintf("%02dm%02d", m, s)
	} else {
		return fmt.Sprintf("%02dh%02dm%02d", h, m, s)
	}
}
