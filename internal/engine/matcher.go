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

func (m *Matcher) Match(line string) bool {
	if m.fixedString {
		if m.ignoreCase {
			return strings.Contains(strings.ToLower(line), m.Pattern)
		}
		return strings.Contains(line, m.Pattern)
	}
	return m.regex.MatchString(line)
}
