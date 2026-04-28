package printer

import (
	"fmt"
	"sort"

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

// outputLine represents a single line to be printed, with enough info
// to decide how to render it (match vs context) and where it sits in the file.
type outputLine struct {
	lineNum int
	content string
	isMatch bool
	ranges  [][2]int
}

func (p *TerminalPrinter) Print(match model.FileMatch) {
	if p.useColor {
		p.fileColor.Println(match.FilePath)
	} else {
		fmt.Println(match.FilePath)
	}

	lines := buildOutputLines(match.Lines)
	blocks := groupIntoBlocks(lines)

	for _, block := range blocks {
		for _, line := range block {
			p.printLine(line)
		}
	}
}

// buildOutputLines converts per-match context into a deduplicated, sorted list of output lines.
// Key rule: if a line is both a context line and a match line, the match wins.
func buildOutputLines(matches []model.LineMatch) []outputLine {
	seen := make(map[int]*outputLine)

	for _, m := range matches {
		// Add before-context lines
		for i, ctx := range m.BeforeCtx {
			ln := m.LineNum - len(m.BeforeCtx) + i
			if _, ok := seen[ln]; !ok {
				seen[ln] = &outputLine{lineNum: ln, content: ctx, isMatch: false}
			}
		}

		// Add match line — match always wins over context
		if existing, ok := seen[m.LineNum]; ok {
			existing.isMatch = true
			existing.content = m.Content
			existing.ranges = m.Ranges
		} else {
			seen[m.LineNum] = &outputLine{lineNum: m.LineNum, content: m.Content, isMatch: true, ranges: m.Ranges}
		}

		// Add after-context lines
		for i, ctx := range m.AfterCtx {
			ln := m.LineNum + 1 + i
			if _, ok := seen[ln]; !ok {
				seen[ln] = &outputLine{lineNum: ln, content: ctx, isMatch: false}
			}
		}
	}

	result := make([]outputLine, 0, len(seen))
	for _, v := range seen {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].lineNum < result[j].lineNum
	})

	return result
}

// groupIntoBlocks splits sorted lines into contiguous blocks.
// If line numbers are not adjacent, a new block starts.
// Between blocks, ripgrep prints "--".
func groupIntoBlocks(lines []outputLine) [][]outputLine {
	if len(lines) == 0 {
		return nil
	}

	var blocks [][]outputLine
	current := []outputLine{lines[0]}

	for i := 1; i < len(lines); i++ {
		if lines[i].lineNum == lines[i-1].lineNum+1 {
			current = append(current, lines[i])
		} else {
			blocks = append(blocks, current)
			current = []outputLine{lines[i]}
		}
	}
	blocks = append(blocks, current)

	return blocks
}

// printLine renders a single line. Match lines use ":" after the line number
// and get highlight; context lines use "-" and print plain text.
// This is the same convention as grep --context.
func (p *TerminalPrinter) printLine(line outputLine) {
	if line.isMatch {
		if p.useColor {
			p.lineNumColor.Printf("%d:", line.lineNum)
		} else {
			fmt.Printf("%d:", line.lineNum)
		}
		p.printHighlighted(line.content, line.ranges)
		fmt.Println()
	} else {
		if p.useColor {
			p.lineNumColor.Printf("%d-", line.lineNum)
		} else {
			fmt.Printf("%d-", line.lineNum)
		}
		fmt.Println(line.content)
	}
}

func (p *TerminalPrinter) printHighlighted(content string, ranges [][2]int) {
	if len(ranges) == 0 {
		fmt.Print(content)
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
