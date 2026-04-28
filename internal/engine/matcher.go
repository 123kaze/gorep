package engine

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher struct {
	Pattern     string
	regex       *regexp.Regexp
	fixedString bool
	ignoreCase  bool
}

func NewMatcher(pattern string, fixedString bool, ignoreCase bool) (*Matcher, error) {
	m := &Matcher{
		Pattern:     pattern,
		fixedString: fixedString,
		ignoreCase:  ignoreCase,
	}

	if fixedString {
		if ignoreCase {
			m.Pattern = strings.ToLower(m.Pattern)
		}
		return m, nil
	}

	regexPattern := pattern
	if ignoreCase {
		regexPattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	m.regex = re // 缓存用于提高下次的对象速度
	return m, nil
}

//func (m *Matcher) Match(line string) bool {
//	if m.fixedString {
//		if m.ignoreCase {
//			return strings.Contains(strings.ToLower(line), m.Pattern)
//		}
//		return strings.Contains(line, m.Pattern)
//	}
//	return m.regex.MatchString(line)
//}

func (m *Matcher) FindHighlightRanges(line string) [][2]int {

	if m.fixedString {
		return findAllFixed(line, m.Pattern, m.ignoreCase)
	}
	locs := m.regex.FindAllStringIndex(line, -1)
	if locs == nil {
		return nil
	}
	ranges := make([][2]int, len(locs))
	for i, loc := range locs {
		ranges[i] = [2]int{loc[0], loc[1]}
	}
	return ranges
}

func findAllFixed(line string, pattern string, ignoreCase bool) [][2]int {
	if pattern == "" {
		return nil
	}
	searchLine := line
	searchPattern := pattern
	if ignoreCase {
		searchLine = strings.ToLower(searchLine)
	}
	var ranges [][2]int
	start := 0
	n := len(searchLine)
	for {
		idx := strings.Index(searchLine[start:], searchPattern)
		if idx == -1 {
			break
		}
		matchStart := start + idx
		matchEnd := matchStart + len(searchPattern)
		ranges = append(ranges, [2]int{matchStart, matchEnd})
		start = matchEnd
		if start >= n {
			break
		}
	}
	return ranges
}
