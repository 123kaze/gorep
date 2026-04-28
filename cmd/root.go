package cmd

import (
	"fmt"
	"gorep/internal/config"
	"gorep/internal/engine"
	"gorep/internal/printer"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	ignoreCase  bool
	fixedString bool
	workers     int
	include     []string
	exclude     []string
	all         bool
	noColor     bool
)
var rootCmd = &cobra.Command{
	Use:   "gorep <pattern> <path>",
	Short: "A fast file search CLI for Go",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.SearchConfig{
			Pattern:     args[0],
			Path:        args[1],
			IgnoreCase:  ignoreCase,
			FixedString: fixedString,
			Workers:     workers,
			Includes:    include,
			Excludes:    exclude,
			All:         all,
			NoColor:     noColor,
		}

		matcher, err := engine.NewMatcher(cfg.Pattern, cfg.FixedString, cfg.IgnoreCase)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		matches, err := engine.SearchPath(cfg, matcher)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if len(matches) == 0 {
			fmt.Println("no matched content was found")
			return
		}
		var p printer.Printer
		p = printer.NewTerminalPrinter(!noColor)
		for _, match := range matches {
			p.Print(match)
		}
		// 初版代码
		//if match == nil {
		//	fmt.Println("no matched content in", cfg.Path)
		//	return
		//}
		//fmt.Println(match.FilePath)
		//for _, line := range match.Lines {
		//	fmt.Printf("  %d:%s\n", line.LineNum, line.Content)
		//}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&ignoreCase, "ignore-case", "i", false, "ignore case distinction")
	rootCmd.Flags().BoolVarP(&fixedString, "fixed-string", "F", false, "target Pattern as a literal string")
	rootCmd.Flags().IntVarP(&workers, "workers", "j", runtime.NumCPU(), "number of worker goroutines")
	rootCmd.Flags().StringSliceVar(&include, "include", nil, "include files matching glob pattern")
	rootCmd.Flags().StringSliceVar(&exclude, "exclude", nil, "exclude files matching glob pattern")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "search for all files")
	rootCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "disable color output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
