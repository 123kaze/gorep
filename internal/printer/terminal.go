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
	useColor     bool
}

func NewTerminalPrinter(useColor bool) *TerminalPrinter {
	p := &TerminalPrinter{
		useColor: useColor,
	}
	if useColor {
		p.fileColor = color.New(color.FgMagenta, color.Bold)
		p.lineNumColor = color.New(color.FgGreen)
		p.matchColor = color.New(color.FgRed, color.Bold)
	}
	return p
}

func (p *TerminalPrinter) Print(match model.FileMatch) {
	if p.useColor {
		p.fileColor.Println(match.FilePath)
	} else {
		fmt.Println(match.FilePath)
	}

	for _, line := range match.Lines {
		if p.useColor {
			p.lineNumColor.Printf(" %d: ", line.LineNum)
		} else {
			fmt.Println(line.LineNum)
		}
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
		if p.useColor {
			p.matchColor.Print(content[r[0]:r[1]])
		} else {
			fmt.Print(content[r[0]:r[1]])
		}
		lastEnd = r[1]
	}
	fmt.Print(content[lastEnd:])
}
