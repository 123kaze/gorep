package filter

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type gitPattern struct {
	raw     string // 原始文本，方便调试
	negate  bool   // 是否 ! 开头
	dirOnly bool   // 是否 / 结尾
	pattern string // 处理后的匹配模式
}

type GitIgnore struct {
	patterns []gitPattern
	root     string
}

func NewGitIgnore(root string) (*GitIgnore, error) {
	path := filepath.Join(root, ".gitignore")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	g := &GitIgnore{
		root: root,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p := gitPattern{}
		p.raw = line
		if strings.HasPrefix(line, "!") {
			p.negate = true
			line = strings.TrimPrefix(line, "!")
		}
		if strings.HasSuffix(line, "/") {
			p.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}
		line = strings.TrimPrefix(line, "/")

		p.pattern = filepath.FromSlash(line)
		g.patterns = append(g.patterns, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GitIgnore) Match(path string, isDir bool) bool {
	rel, err := filepath.Rel(g.root, path)
	if err != nil {
		return false
	}
	rel = filepath.Clean(rel)

	ignored := false

	for _, p := range g.patterns {
		if p.dirOnly && !isDir {
			continue
		}
		if matchGitIgnorePattern(rel, p.pattern) {
			ignored = !p.negate
		}
	}
	return ignored
}

func matchGitIgnorePattern(rel string, pattern string) bool {
	if pattern == "" {
		return false
	}
	if matched, err := filepath.Match(pattern, filepath.Base(rel)); err == nil && matched {
		return true
	}
	if matched, err := filepath.Match(pattern, rel); err == nil && matched {
		return true
	}
	if strings.HasPrefix(rel, pattern+string(filepath.Separator)) {
		return true
	}
	if strings.Contains(pattern, string(filepath.Separator)) {
		return false
	}
	for {
		dir := filepath.Dir(rel)
		if dir == "." || dir == string(filepath.Separator) {
			break
		}
		base := filepath.Base(dir)
		if matched, err := filepath.Match(pattern, base); err == nil && matched {
			return true
		}
		rel = dir
	}
	return false
}
