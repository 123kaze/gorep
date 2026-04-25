package engine

import (
	"bufio"
	"gorep/internal/model"
	"os"
)

func SearchFile(filePath string, macher Matcher) (*model.FileMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // 防止忘了关闭

	var matches []model.LineMatch
	scanner := bufio.NewScanner(file)
	LineNum := 0

	for scanner.Scan() {
		LineNum++
		line := scanner.Text()
		if macher.Match(line) {
			matches = append(matches, model.LineMatch{
				LineNum: LineNum,
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
