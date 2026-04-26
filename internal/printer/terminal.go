package printer

import (
	"fmt"
	"gorep/internal/model"

	"github.com/fatih/color"
)

type TerminalPrinter struct {
	fileColor    *color.Color
	lineNumColor *color.Color
	matchColor   *color.Color
}

func NewTerminalPrinter() *TerminalPrinter {
	return &TerminalPrinter{
		fileColor:    color.New(color.FgMagenta, color.Bold),
		lineNumColor: color.New(color.FgGreen),
		matchColor:   color.New(color.FgRed, color.Bold),
	}
}

func (p *TerminalPrinter) Print(match model.FileMatch) {
	p.fileColor.Println(match.FilePath)
	for _, line := range match.Lines {
		p.lineNumColor.Printf(" %d: ", line.LineNum)
		p.printHighlighted(line.Content, line.Ranges)
		fmt.Println()
	}
}

func (p *TerminalPrinter) printHighlighted(content string, ranges [][2]int) {
	if len(ranges) == 0 {
		//fmt.Print(content)
		return
	}
	lastEnd := 0
	for _, r := range ranges {
		fmt.Print(content[lastEnd:r[0]])
		p.matchColor.Print(content[r[0]:r[1]])
		lastEnd = r[1]
	}
	fmt.Print(content[lastEnd:])
}
