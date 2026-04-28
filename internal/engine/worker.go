package engine

import (
	"bufio"
	"io"
	"os"

	"gorep/internal/model"
)

func SearchFile(filePath string, matcher *Matcher, beforeCtx int, afterCtx int) (*model.FileMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read first 512 bytes for binary detection — only open the file once
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	if containsNullByte(buf[:n]) {
		return nil, nil
	}
	// Seek back to start so the scanner reads from the beginning
	file.Seek(0, io.SeekStart)

	var matches []model.LineMatch
	scanner := bufio.NewScanner(file)
	lineNum := 0
	var beforeBuf []string
	pendingAfter := make(map[int]int)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for idx, remaining := range pendingAfter {
			matches[idx].AfterCtx = append(matches[idx].AfterCtx, line)
			remaining--
			if remaining == 0 {
				delete(pendingAfter, idx)
			} else {
				pendingAfter[idx] = remaining
			}
		}

		ranges := matcher.FindHighlightRanges(line)
		if len(ranges) > 0 {
			matches = append(matches, model.LineMatch{
				LineNum:   lineNum,
				Content:   line,
				Ranges:    ranges,
				BeforeCtx: copyStrings(beforeBuf),
			})
			if afterCtx > 0 {
				pendingAfter[len(matches)-1] = afterCtx
			}
		}
		if beforeCtx > 0 {
			if len(beforeBuf) >= beforeCtx {
				beforeBuf = beforeBuf[1:]
			}
			beforeBuf = append(beforeBuf, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, nil
	}

	return &model.FileMatch{
		FilePath: filePath,
		Lines:    matches,
	}, nil
}

func containsNullByte(buf []byte) bool {
	for _, b := range buf {
		if b == 0 {
			return true
		}
	}
	return false
}

func copyStrings(buf []string) []string {
	if len(buf) == 0 {
		return nil
	}
	dst := make([]string, len(buf))
	copy(dst, buf)
	return dst
}
