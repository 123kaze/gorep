package filter

import (
	"path/filepath"
	"strings"
)

type GlobFilter struct {
	includes []string
	excludes []string
}

func NewGlobFilter(includes, excludes []string) *GlobFilter {
	return &GlobFilter{
		includes: normalizePatterns(includes),
		excludes: normalizePatterns(excludes),
	}
}

func normalizePatterns(patterns []string) []string {
	out := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		out = append(out, filepath.FromSlash(pattern))
	}
	return out
}

func (f *GlobFilter) Match(path string) bool {
	//normalizedPath := filepath.ToSlash(path)
	if len(f.includes) > 0 && !matchAny(path, f.includes) {
		return false
	}
	if matchAny(path, f.excludes) {
		return false
	}
	return true
}

func matchAny(path string, patterns []string) bool {
	for _, p := range patterns {
		if matchGlob(path, p) {
			return true
		}
	}
	return false
}

func matchGlob(path, pattern string) bool {
	matched, err := filepath.Match(pattern, filepath.Base(path))
	if err == nil && matched {
		return true
	}
	matched, err = filepath.Match(pattern, path)
	if err == nil && matched {
		return true
	}
	sep := string(filepath.Separator)
	if strings.HasSuffix(pattern, sep+"*") {
		prefix := strings.TrimSuffix(pattern, sep+"*")
		return strings.HasPrefix(path, prefix+sep)
	}
	return false
}
