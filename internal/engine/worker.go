package engine

import (
	"bufio"
	"os"

	"gorep/internal/model"
)

func SearchFile(filePath string, matcher *Matcher) (*model.FileMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matches []model.LineMatch
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if matcher.Match(line) {
			matches = append(matches, model.LineMatch{
				LineNum: lineNum,
				Content: line,
			})
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
