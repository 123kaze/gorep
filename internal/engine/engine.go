package engine

import (
	"fmt"
	"gorep/internal/config"
	"gorep/internal/filter"
	"gorep/internal/model"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func SearchPath(cfg config.SearchConfig, matcher *Matcher) ([]model.FileMatch, error) {
	info, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		match, err := SearchFile(cfg.Path, matcher, cfg.BeforeCtx, cfg.AfterCtx)
		if err != nil {
			return nil, err
		}
		if match == nil {
			return nil, nil
		}
		return []model.FileMatch{*match}, nil
	}
	// 新增用cfg传递参数，然后用chan来做多线程
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	fileCh := make(chan string, cfg.Workers*2)
	resultCh := make(chan *model.FileMatch, cfg.Workers)
	globFilter := filter.NewGlobFilter(cfg.Includes, cfg.Excludes)
	var gitIgnore *filter.GitIgnore
	if !cfg.NoGitIgnore {
		gitIgnore, _ = filter.NewGitIgnore(cfg.Path)
	}
	var wg sync.WaitGroup

	go func() {
		defer close(fileCh)
		walkFiles(cfg.Path, globFilter, gitIgnore, fileCh, cfg.All)
	}()

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for filepath := range fileCh {
				match, err := SearchFile(filepath, matcher, cfg.BeforeCtx, cfg.AfterCtx)
				if err != nil {
					fmt.Printf("Error searching %s: %s\n", filepath, err)
					continue
				}
				if match != nil {
					resultCh <- match
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []model.FileMatch
	for match := range resultCh {
		results = append(results, *match)
	}

	return results, nil
	// version 1.0 代码
	//var results []model.FileMatch
	//err = filepath.WalkDir(cfg.Path, func(currentPath string, d fs.DirEntry, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	if d.IsDir() {
	//		return nil
	//	}
	//	match, err := SearchFile(currentPath, matcher)
	//	if err != nil {
	//		fmt.Println("error:", err)
	//		return err
	//	}
	//	if match == nil {
	//		return nil
	//	}
	//	results = append(results, *match)
	//	return nil
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return results, nil
}

func walkFiles(root string, globFilter *filter.GlobFilter, gitIgnore *filter.GitIgnore, fileCh chan<- string, all bool) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			isHidden := path != root && strings.HasPrefix(d.Name(), ".")
			if !all && isHidden {
				return filepath.SkipDir
			}
			if gitIgnore != nil && gitIgnore.Match(path, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if !all && strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if gitIgnore != nil && gitIgnore.Match(path, false) {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		if !globFilter.Match(relPath) {
			return nil
		}
		if filter.HasBinaryExt(path) {
			return nil
		}
		fileCh <- path
		return nil
	})
}
